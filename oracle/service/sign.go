package service

import (
	"context"
	"errors"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gofiber/fiber/v2/log"
	"github.com/holiman/uint256"
	"github.com/near/borsh-go"
	blockchain2 "gringotts/blockchain"
	"gringotts/connection"
	"gringotts/models"
	"gringotts/provider"
	"gringotts/utils"
	"math/big"
	"strings"
)

func CreateBlockchainTransaction(amount *big.Int, servingContext *models.ServingContext) error {
	var err error
	var tx *models.UnsignedTransaction

	switch servingContext.SourceChain {
	case models.Solana, models.SolanaDev:
		tx, err = CreateSolanaTransaction(amount, servingContext)
	default:
		tx, err = CreateEVMTransaction(amount, servingContext)
	}

	if err != nil {
		return err
	}

	servingContext.Tx = tx
	return nil
}

func CreateSolanaTransaction(
	amount *big.Int,
	servingContext *models.ServingContext,
) (*models.UnsignedTransaction, error) {
	inboundTransfers := make([]connection.BridgeInboundTransferItem, len(servingContext.Inbounds))
	for i, inbound := range servingContext.Inbounds {
		inboundTransfer := connection.BridgeInboundTransferItem{
			Asset:  utils.ToByte32SOL(inbound.Token.Address),
			Amount: inbound.Amount.Uint64(),
		}

		if inbound.Swap != nil {
			inboundTransfer.Swap = &connection.Swap{
				Executor:      utils.ToByte32SOL(inbound.Swap.Executor),
				Command:       utils.FromHex(inbound.Swap.Command),
				AccountsCount: uint8(len(inbound.Swap.Accounts)),
				StableToken:   utils.ToByte32SOL(inbound.Swap.ToToken.Address),
			}
		}

		inboundTransfers[i] = inboundTransfer
	}

	outboundTransfers := make([]connection.BridgeOutboundTransfer, 0, len(servingContext.Outbounds))
	for chain, transactions := range servingContext.Outbounds {
		outboundTransferItems := make([]connection.BridgeOutboundTransferItem, len(transactions))
		totalGas := uint64(0)

		for i, item := range transactions {
			gas, _ := GetExecutionParams(chain, item.Token)
			totalGas = totalGas + gas

			outboundTransferItem := connection.BridgeOutboundTransferItem{
				DistributionBP: uint16(item.DistributionBPS),
			}

			outboundTransferItems[i] = outboundTransferItem
		}

		outboundTransfers = append(outboundTransfers, connection.BridgeOutboundTransfer{
			ChainID:      chain.GetId(),
			ExecutionGas: totalGas,
			Message:      CreateMessage(chain, transactions),
			Items:        outboundTransferItems,
		})
	}

	requestSerializedData, _ := borsh.Serialize(connection.BridgeRequest{
		Inbound: connection.BridgeInboundTransfer{
			AmountUSDx: amount.Uint64(),
			Items:      inboundTransfers,
		},
		Outbounds: outboundTransfers,
	})
	data := append(models.GetBridgeDiscriminator(), requestSerializedData...)

	solanaStableCoin := models.GetDefaultStableCoins(servingContext.SourceChain)

	accounts := []types.AccountMeta{
		{PubKey: solana.PublicKeyFromString(servingContext.Wallet), IsSigner: true, IsWritable: true},                                                                                        // user
		{PubKey: solana.PublicKeyFromString(models.GetPriceFeed()), IsSigner: false, IsWritable: false},                                                                                      // pricefeed
		{PubKey: solana.PublicKeyFromString(models.GetGringotts(servingContext.SourceChain)), IsSigner: false, IsWritable: false},                                                            // gringotts
		{PubKey: solana.PublicKeyFromString(models.GetVault(servingContext.SourceChain)), IsSigner: false, IsWritable: true},                                                                 // vault
		{PubKey: solana.PublicKeyFromString(models.GetPeer(servingContext.SourceChain, servingContext.SourceChain)), IsSigner: false, IsWritable: false},                                     // self
		{PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(models.GetGringotts(servingContext.SourceChain), solanaStableCoin.Address)), IsSigner: false, IsWritable: true}, // gringotts_stable_coin
		{PubKey: solana.PublicKeyFromString(solanaStableCoin.Address), IsSigner: false, IsWritable: false},                                                                                   // mint
		{PubKey: solana.PublicKeyFromString(models.GetJupiter()), IsSigner: false, IsWritable: false},                                                                                        // swap
		{PubKey: solana.SPLAssociatedTokenAccountProgramID, IsSigner: false, IsWritable: false},                                                                                              // spl
		{PubKey: solana.TokenProgramID, IsSigner: false, IsWritable: false},                                                                                                                  // token
		{PubKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},                                                                                                                 // system
	}

	alts := make([]string, 0)
	for _, item := range servingContext.Inbounds {
		if item.Token.Address == solanaStableCoin.Address {
			accounts = append(accounts, types.AccountMeta{
				PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(servingContext.Wallet, solanaStableCoin.Address)), IsSigner: false, IsWritable: true,
			})
		} else {
			if item.Token.IsStableCoin {
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.NativeMint), IsSigner: false, IsWritable: false,
				})
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(models.GetGringotts(servingContext.SourceChain), models.NativeMint)), IsSigner: false, IsWritable: true,
				})
			} else {
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(item.Token.Address), IsSigner: false, IsWritable: false,
				})
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(models.GetGringotts(servingContext.SourceChain), item.Token.Address)), IsSigner: false, IsWritable: true,
				})
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(servingContext.Wallet, item.Token.Address)), IsSigner: false, IsWritable: true,
				})
			}

			for _, acc := range item.Swap.Accounts {
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(acc.Address), IsSigner: acc.IsSigner, IsWritable: acc.IsWritable,
				})
			}
			alts = append(alts, item.Swap.AddressLookup...)
		}
	}
	for chain := range servingContext.Outbounds {
		for _, acc := range models.GetSendAccounts(chain) {
			accounts = append(accounts, types.AccountMeta{
				PubKey: solana.PublicKeyFromString(acc.Address), IsSigner: acc.IsSigner, IsWritable: acc.IsWritable,
			})
		}
	}

	instruction := types.Instruction{
		ProgramID: solana.PublicKeyFromString(servingContext.SourceChain.GetContract()),
		Accounts:  accounts,
		Data:      data,
	}

	client := blockchain2.GetSOLConnection(servingContext.SourceChain)
	recentBlockhash, _ := client.GetLatestBlockhash(context.Background())

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:                   solana.PublicKeyFromString(servingContext.Wallet),
		RecentBlockhash:            recentBlockhash.Blockhash,
		Instructions:               []types.Instruction{instruction},
		AddressLookupTableAccounts: blockchain2.GetALT(servingContext.SourceChain, alts),
	})

	tx, _ := message.Serialize()

	return &models.UnsignedTransaction{
		Data: tx,
	}, nil
}

func CreateEVMTransaction(amount *big.Int, servingContext *models.ServingContext) (*models.UnsignedTransaction, error) {
	value := uint256.NewInt(0)

	inboundTransfers := make([]connection.GringottsBridgeInboundTransferItem, len(servingContext.Inbounds))

	for i, inbound := range servingContext.Inbounds {
		transactionItem := connection.GringottsBridgeInboundTransferItem{
			Asset:  utils.ToByte32(inbound.Token.Address),
			Amount: inbound.Amount.ToBig(),
		}

		if inbound.Swap != nil {
			transactionItem.Command = utils.FromHex(inbound.Swap.Command)
			transactionItem.Executor = utils.ToByte32(inbound.Swap.Executor)
			transactionItem.StableToken = utils.ToByte32(inbound.Swap.ToToken.Address)
		}

		if inbound.Token.IsNative {
			value = inbound.Amount
		}

		inboundTransfers[i] = transactionItem
	}

	outboundTransfers := make([]connection.GringottsBridgeOutboundTransfer, 0, len(servingContext.Outbounds))
	for chain, items := range servingContext.Outbounds {
		transactionItems := make([]connection.GringottsBridgeOutboundTransferItem, len(items))
		totalGas := uint64(0)

		for i, item := range items {
			gas, _ := GetExecutionParams(chain, item.Token)
			totalGas = totalGas + gas

			transactionItems[i] = connection.GringottsBridgeOutboundTransferItem{
				DistributionBP: uint16(item.DistributionBPS),
			}
		}

		outboundTransfers = append(outboundTransfers, connection.GringottsBridgeOutboundTransfer{
			ChainId:      chain.GetId(),
			Message:      CreateMessage(chain, items),
			Items:        transactionItems,
			ExecutionGas: big.NewInt(int64(totalGas)),
		})
	}

	request := connection.GringottsBridgeRequest{
		Inbound: connection.GringottsBridgeInboundTransfer{
			AmountUSDX: amount,
			Items:      inboundTransfers,
		},
		Outbounds: outboundTransfers,
	}

	parsedABI, _ := abi.JSON(strings.NewReader(connection.GringottsEVMMetaData.ABI))
	data, _ := parsedABI.Pack("bridge", request)

	return &models.UnsignedTransaction{
		Data:  data,
		Value: value,
	}, nil
}

func CreateMessage(chain models.Blockchain, items []*models.Outbound) []byte {
	switch chain {
	case models.Solana, models.SolanaDev:
		message, err := CreateSolanaMessage(chain, items)
		if err != nil {
			log.Errorw("failed to create solana message", "error", err)
			return nil
		}
		return message
	default:
		return CreateEVMMessage(items)
	}
}

func CreateEVMMessage(items []*models.Outbound) []byte {
	message := make([]byte, 0, 100)
	message = append(message, byte(len(items)))

	for _, item := range items {
		message = append(message, utils.FromByte32ToByte(item.Token.Address)...)
		message = append(message, utils.FromByte32ToByte(item.Recipient)...)

		if item.Swap != nil {
			message = append(message, utils.FromByte32ToByte(item.Swap.Executor)...)
			message = append(message, utils.FromByte32ToByte(item.Swap.FromToken.Address)...)
			message = append(message, utils.FromByte32ToByte(item.Swap.Command)...)
		}
	}

	return message
}

func CreateSolanaMessage(chain models.Blockchain, items []*models.Outbound) ([]byte, error) {
	var mainFromToken *models.Token

	for _, item := range items {
		var curMainFromToken *models.Token

		if item.Swap == nil {
			curMainFromToken = item.Token
		} else {
			curMainFromToken = item.Swap.FromToken
		}

		if mainFromToken == nil {
			mainFromToken = curMainFromToken
		} else {
			if mainFromToken.Address != curMainFromToken.Address {
				return nil, errors.New("not single stable coin")
			}
		}
	}

	if mainFromToken == nil {
		return nil, errors.New("no stable coin")
	}

	stableIndex := 0
	for i, stableCoin := range models.GetStableCoins(chain) {
		if stableCoin.Address == mainFromToken.Address {
			stableIndex = i
			break
		}
	}

	metadata := []byte{byte(stableIndex)}

	allAccounts := make([]*models.Account, 0)
	for _, transaction := range items {
		allAccounts = append(allAccounts,
			&models.Account{Address: transaction.Recipient, IsWritable: transaction.Token.IsNative},
		)

		if transaction.Swap == nil {
			allAccounts = append(allAccounts,
				&models.Account{Address: mainFromToken.Address},
				&models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, mainFromToken.Address)},
			)
		} else {
			if transaction.Token.IsNative {
				allAccounts = append(allAccounts,
					&models.Account{Address: models.NativeMint},
					&models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint), IsWritable: true},
				)
			} else {
				allAccounts = append(allAccounts,
					&models.Account{Address: transaction.Token.Address},
					&models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.Token.Address), IsWritable: true},
				)
			}

			for _, acc := range transaction.Swap.Accounts {
				allAccounts = append(allAccounts,
					&models.Account{Address: acc.Address, IsWritable: acc.IsWritable},
				)
			}
		}
	}

	gringottsStable := models.GetAssociatedTokenAddress(models.GetGringotts(chain), mainFromToken.Address)

	dontSendAccount := map[string]bool{
		models.GetGringotts(chain):                         true,
		models.GetVault(chain):                             true,
		solana.SPLAssociatedTokenAccountProgramID.String(): true,
		solana.TokenProgramID.String():                     true,
		solana.SystemProgramID.String():                    true,
		mainFromToken.Address:                              true,
		gringottsStable:                                    true,
		provider.JupiterAddress:                            true,
	}

	for _, transaction := range items {
		if transaction.Swap == nil {
			userStable := models.GetAssociatedTokenAddress(transaction.Recipient, mainFromToken.Address)
			dontSendAccount[userStable] = true
		} else {
			if transaction.Token.IsNative {
				gringottsWSOL := models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint)
				dontSendAccount[gringottsWSOL] = true

			} else {
				associatedToken := models.GetAssociatedTokenAddress(transaction.Recipient, transaction.Token.Address)
				dontSendAccount[associatedToken] = true
			}
		}
	}

	// Accounts need to be sent
	sendingAccountsMap := make(map[string]*models.Account)
	for _, account := range allAccounts {
		if dontSendAccount[account.Address] {
			continue
		}

		if _, ok := sendingAccountsMap[account.Address]; ok {
			sendingAccountsMap[account.Address].IsWritable = sendingAccountsMap[account.Address].IsWritable || account.IsWritable
		} else {
			sendingAccountsMap[account.Address] = account
		}
	}

	addressMap := map[string]byte{
		models.GetGringotts(chain):                         0,
		solana.SPLAssociatedTokenAccountProgramID.String(): 1,
		solana.TokenProgramID.String():                     2,
		solana.SystemProgramID.String():                    3,
		provider.JupiterAddress:                            4,
		mainFromToken.Address:                              5,
		gringottsStable:                                    6,
	}

	index := byte(len(addressMap))

	flags := ""
	metadata = append(metadata, byte(len(sendingAccountsMap)))
	for _, account := range sendingAccountsMap {
		metadata = append(metadata, utils.FromByte32ToByteSOL(account.Address)...)

		addressMap[account.Address] = index
		index++

		if account.IsWritable {
			flags = flags + "1"
		} else {
			flags = flags + "0"
		}
	}
	metadata = append(metadata, utils.ZeroOneStringToByteArray(flags)...)

	for _, transaction := range items {
		address := ""

		if transaction.Swap == nil {
			address = models.GetAssociatedTokenAddress(transaction.Recipient, mainFromToken.Address)
		} else {
			if transaction.Token.IsNative {
				address = models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint)
			} else {
				address = models.GetAssociatedTokenAddress(transaction.Recipient, transaction.Token.Address)
			}
		}

		if _, ok := addressMap[address]; !ok && address != "" {
			addressMap[address] = index
			index++
		}
	}

	for _, transaction := range items {
		metadata = append(metadata, addressMap[transaction.Recipient])

		if transaction.Swap == nil {
			userStable := models.GetAssociatedTokenAddress(transaction.Recipient, mainFromToken.Address)

			metadata = append(metadata, addressMap[mainFromToken.Address])
			metadata = append(metadata, addressMap[userStable])
		} else {
			if transaction.Token.IsNative {
				gringottsWSOL := models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint)

				metadata = append(metadata, addressMap[models.NativeMint])
				metadata = append(metadata, addressMap[gringottsWSOL])
			} else {
				associatedToken := models.GetAssociatedTokenAddress(transaction.Recipient, transaction.Token.Address)

				metadata = append(metadata, addressMap[transaction.Token.Address])
				metadata = append(metadata, addressMap[associatedToken])
			}

			metadata = append(metadata, byte(len(transaction.Swap.Accounts)))
			for _, account := range transaction.Swap.Accounts {
				metadata = append(metadata, addressMap[account.Address])
			}

			command := utils.FromHex(transaction.Swap.Command)
			metadata = append(metadata, utils.ToBigEndianBytes(uint16(len(command)))...)
			metadata = append(metadata, command...)
		}
	}

	return metadata, nil
}
