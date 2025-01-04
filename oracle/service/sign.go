package service

import (
	"context"
	"fmt"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func CreateBlockchainTransaction(
	amount *big.Int,
	wallet string,
	blockchain models.Blockchain,
	inTransactions map[models.Blockchain][]*models.Transaction,
	outTransactions map[models.Blockchain][]*models.Transaction,
) (*provider.UnsignedTransaction, error) {
	switch blockchain {
	case models.Solana, models.SolanaDev:
		return createSolanaTransaction(amount, wallet, blockchain, inTransactions, outTransactions)
	default:
		return createEVMTransaction(amount, blockchain, inTransactions, outTransactions)
	}
}

func createSolanaTransaction(
	amount *big.Int,
	wallet string,
	chain models.Blockchain,
	inTransactions map[models.Blockchain][]*models.Transaction,
	outTransactions map[models.Blockchain][]*models.Transaction,
) (*provider.UnsignedTransaction, error) {
	inTransaction := inTransactions[chain]
	inboundTransfers := make([]connection.BridgeInboundTransferItem, len(inTransaction))
	for i, transaction := range inTransaction {
		tAmount, _ := uint256.FromDecimal(transaction.SrcAmount)
		inboundTransfer := connection.BridgeInboundTransferItem{
			Asset:  utils.ToByte32SOL(transaction.FromToken),
			Amount: tAmount.Uint64(),
		}

		if transaction.Swap != nil {
			inboundTransfer.Swap = &connection.Swap{
				Executor:      utils.ToByte32SOL(transaction.Swap.Executor),
				Command:       utils.FromHex(transaction.Swap.Command),
				AccountsCount: uint8(len(transaction.Swap.Accounts)),
				StableToken:   utils.ToByte32SOL(transaction.ToToken),
			}
		}

		inboundTransfers[i] = inboundTransfer
	}

	outboundTransfers := make([]connection.BridgeOutboundTransfer, 0, len(outTransactions))
	for chainIter, transactions := range outTransactions {
		outboundTransferItems := make([]connection.BridgeOutboundTransferItem, len(transactions))
		totalGas := uint64(0)

		for i, transaction := range transactions {
			gas, _, _ := GetExecutionParams(chainIter, transaction.ToToken)

			outboundTransferItem := connection.BridgeOutboundTransferItem{
				DistributionBP: uint16(transaction.DistributionBPS),
			}

			outboundTransferItems[i] = outboundTransferItem
			totalGas = totalGas + gas
		}

		outboundTransfers = append(outboundTransfers, connection.BridgeOutboundTransfer{
			ChainID:      chainIter.GetId(),
			ExecutionGas: totalGas,
			Message:      getMessage(chainIter, transactions),
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

	solanaStableCoin := models.GetDefaultStableCoins(chain)

	accounts := []types.AccountMeta{
		{PubKey: solana.PublicKeyFromString(wallet), IsSigner: true, IsWritable: true},                                                                                  // user
		{PubKey: solana.PublicKeyFromString(models.GetPriceFeed()), IsSigner: false, IsWritable: false},                                                                 // pricefeed
		{PubKey: solana.PublicKeyFromString(models.GetGringotts(chain)), IsSigner: false, IsWritable: false},                                                            // gringotts
		{PubKey: solana.PublicKeyFromString(models.GetVault(chain)), IsSigner: false, IsWritable: true},                                                                 // vault
		{PubKey: solana.PublicKeyFromString(models.GetPeer(chain, chain)), IsSigner: false, IsWritable: false},                                                          // self
		{PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(models.GetGringotts(chain), solanaStableCoin.Address)), IsSigner: false, IsWritable: true}, // gringotts_stable_coin
		{PubKey: solana.PublicKeyFromString(solanaStableCoin.Address), IsSigner: false, IsWritable: false},                                                              // mint
		{PubKey: solana.PublicKeyFromString(models.GetJupiter()), IsSigner: false, IsWritable: false},                                                                   // swap
		{PubKey: solana.SPLAssociatedTokenAccountProgramID, IsSigner: false, IsWritable: false},                                                                         // spl
		{PubKey: solana.TokenProgramID, IsSigner: false, IsWritable: false},                                                                                             // token
		{PubKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},                                                                                            // system
	}

	alts := make([]string, 0)
	for _, transaction := range inTransaction {
		inToken := models.GetToken(chain, transaction.FromToken)

		if inToken.Address == solanaStableCoin.Address {
			accounts = append(accounts, types.AccountMeta{
				PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(wallet, solanaStableCoin.Address)), IsSigner: false, IsWritable: true,
			})
		} else {
			if inToken.Address == "" {
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.NativeMint), IsSigner: false, IsWritable: false,
				})
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint)), IsSigner: false, IsWritable: true,
				})
			} else {
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(inToken.Address), IsSigner: false, IsWritable: false,
				})
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(models.GetGringotts(chain), inToken.Address)), IsSigner: false, IsWritable: true,
				})
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(models.GetAssociatedTokenAddress(wallet, inToken.Address)), IsSigner: false, IsWritable: true,
				})
			}

			for _, acc := range transaction.Swap.Accounts {
				accounts = append(accounts, types.AccountMeta{
					PubKey: solana.PublicKeyFromString(acc.Address), IsSigner: acc.IsSigner, IsWritable: acc.IsWritable,
				})
			}
			alts = append(alts, transaction.Swap.AddressLookup...)
		}
	}
	for chainIter := range outTransactions {
		for _, acc := range models.GetSendAccounts(chainIter) {
			accounts = append(accounts, types.AccountMeta{
				PubKey: solana.PublicKeyFromString(acc.Address), IsSigner: acc.IsSigner, IsWritable: acc.IsWritable,
			})
		}
	}

	instruction := types.Instruction{
		ProgramID: solana.PublicKeyFromString(chain.GetContract()),
		Accounts:  accounts,
		Data:      data,
	}

	client := blockchain2.GetSOLConnection(chain)
	recentBlockhash, _ := client.GetLatestBlockhash(context.Background())

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:                   solana.PublicKeyFromString(wallet),
		RecentBlockhash:            recentBlockhash.Blockhash,
		Instructions:               []types.Instruction{instruction},
		AddressLookupTableAccounts: blockchain2.GetALT(chain, alts),
	})

	tx, _ := message.Serialize()

	return &provider.UnsignedTransaction{
		Contract: chain.GetContract(),
		Data:     tx,
	}, nil
}

func createEVMTransaction(
	amount *big.Int,
	blockchain models.Blockchain,
	inTransactions map[models.Blockchain][]*models.Transaction,
	outTransactions map[models.Blockchain][]*models.Transaction,
) (*provider.UnsignedTransaction, error) {
	inTransaction := inTransactions[blockchain]
	value := ""

	inboundTransfers := make([]connection.GringottsBridgeInboundTransferItem, len(inTransaction))
	for i, transaction := range inTransaction {
		tAmount, _ := uint256.FromDecimal(transaction.SrcAmount)
		transactionItem := connection.GringottsBridgeInboundTransferItem{
			Asset:  utils.ToByte32(transaction.FromToken),
			Amount: tAmount.ToBig(),
		}

		if transaction.Swap != nil {
			transactionItem.Command = utils.FromHex(transaction.Swap.Command)
			transactionItem.Executor = utils.ToByte32(transaction.Swap.Executor)
			transactionItem.StableToken = utils.ToByte32(transaction.ToToken)
		}

		if len(transaction.FromToken) == 0 {
			value = transaction.SrcAmount
		}

		inboundTransfers[i] = transactionItem
	}

	outboundTransfers := make([]connection.GringottsBridgeOutboundTransfer, 0, len(outTransactions))
	for chain, transactions := range outTransactions {
		transactionItems := make([]connection.GringottsBridgeOutboundTransferItem, len(transactions))
		totalGas := int64(0)

		for i, transaction := range transactions {
			gas, _, _ := GetExecutionParams(chain, transaction.ToToken)
			totalGas = totalGas + int64(gas)

			transactionItems[i] = connection.GringottsBridgeOutboundTransferItem{
				DistributionBP: uint16(transaction.DistributionBPS),
			}
		}

		outboundTransfers = append(outboundTransfers, connection.GringottsBridgeOutboundTransfer{
			ChainId:      chain.GetId(),
			Message:      getMessage(chain, transactions),
			Items:        transactionItems,
			ExecutionGas: big.NewInt(totalGas),
		})

		//fmt.Println("Message -> ", utils.ToHex(getMessage(chain, transactions)))
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

	return &provider.UnsignedTransaction{
		Contract: blockchain.GetContract(),
		Data:     data,
		Value:    value,
	}, nil
}

func getMessage(chain models.Blockchain, transactions []*models.Transaction) []byte {
	switch chain {
	case models.Solana, models.SolanaDev:
		return getMetaData(chain, transactions)
	}

	message := make([]byte, 0, 100)
	message = append(message, byte(len(transactions)))

	for _, transaction := range transactions {
		message = append(message, utils.FromByte32ToByte(transaction.ToToken)...)
		message = append(message, utils.FromByte32ToByte(transaction.Recipient)...)

		if transaction.Swap == nil {
			message = append(message, 0)
		} else {
			message = append(message, 1)
			message = append(message, utils.FromByte32ToByte(transaction.Swap.Executor)...)
			message = append(message, utils.FromByte32ToByte(transaction.FromToken)...)
			message = append(message, utils.FromByte32ToByte(transaction.Swap.Command)...)
		}
	}

	return message
}

func getMetaData(chain models.Blockchain, transactions []*models.Transaction) []byte {
	accounts := make([]*models.Account, 0)

	for _, transaction := range transactions {
		accounts = append(accounts,
			&models.Account{Address: transaction.Recipient, IsSigner: false, IsWritable: transaction.ToToken == ""},
		)

		if transaction.Swap == nil {
			accounts = append(accounts,
				&models.Account{Address: transaction.ToToken, IsSigner: false, IsWritable: false},

				&models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), transaction.ToToken), IsSigner: false, IsWritable: true},
				&models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.ToToken), IsSigner: false, IsWritable: true},
			)
		} else {
			if transaction.ToToken != "" {
				accounts = append(accounts,
					&models.Account{Address: transaction.ToToken, IsSigner: false, IsWritable: false},
					&models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.ToToken), IsSigner: false, IsWritable: true},
				)
			} else {
				accounts = append(accounts,
					&models.Account{Address: models.NativeMint, IsSigner: false, IsWritable: false},
					&models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint), IsSigner: false, IsWritable: true},
				)
			}

			accounts = append(accounts,
				&models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.FromToken), IsSigner: false, IsWritable: true},

				&models.Account{Address: transaction.FromToken, IsSigner: false, IsWritable: false},
				&models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), transaction.FromToken), IsSigner: false, IsWritable: true},

				&models.Account{Address: provider.JupiterAddress, IsSigner: false, IsWritable: false},
			)

			for _, acc := range transaction.Swap.Accounts {
				accounts = append(accounts,
					&models.Account{Address: acc.Address, IsSigner: false, IsWritable: acc.IsWritable},
				)
			}
		}
	}

	dontSendAccount := map[string]*models.Account{
		models.GetGringotts(chain):                         {Address: models.GetGringotts(chain)},
		models.GetVault(chain):                             {Address: models.GetVault(chain)},
		solana.SPLAssociatedTokenAccountProgramID.String(): {Address: solana.SPLAssociatedTokenAccountProgramID.String()},
		solana.TokenProgramID.String():                     {Address: solana.TokenProgramID.String()},
		solana.SystemProgramID.String():                    {Address: solana.SystemProgramID.String()},
	}

	accountsMap := make(map[string]*models.Account)
	for _, account := range accounts {
		if _, ok := dontSendAccount[account.Address]; ok {
			continue
		}

		if _, ok := accountsMap[account.Address]; ok {
			accountsMap[account.Address].IsWritable = accountsMap[account.Address].IsWritable || account.IsWritable
		} else {
			accountsMap[account.Address] = account
		}
	}

	addressMap := map[string]int{
		models.GetGringotts(chain):                         0,
		models.GetVault(chain):                             1,
		solana.SPLAssociatedTokenAccountProgramID.String(): 2,
		solana.TokenProgramID.String():                     3,
		solana.SystemProgramID.String():                    4,
	}

	index := 5
	flags := ""

	metadata := []byte{byte(len(accountsMap))}
	for _, account := range accountsMap {
		metadata = append(metadata, utils.FromByte32ToByteSOL(account.Address)...)

		addressMap[account.Address] = index
		index++

		if account.IsWritable {
			flags = flags + "1"
		} else {
			flags = flags + "0"
		}

		fmt.Println(account.Address)
	}
	fmt.Println(len(accountsMap))

	metadata = append(metadata, utils.ZeroOneStringToByteArray(flags)...)
	metadata = append(metadata, byte(len(accounts)))

	for _, account := range accounts {
		metadata = append(metadata, byte(addressMap[account.Address]))
	}

	itemsFlag := ""
	for _, transaction := range transactions {
		if transaction.Swap == nil {
			itemsFlag = itemsFlag + "0"
		} else {
			itemsFlag = itemsFlag + "1"
		}
	}

	metadata = append(metadata, utils.ZeroOneStringToByteArray(itemsFlag)...)

	for _, transaction := range transactions {
		if transaction.Swap == nil {
			continue
		} else {
			metadata = append(metadata, byte(len(transaction.Swap.Accounts)))

			command := utils.FromHex(transaction.Swap.Command)
			metadata = append(metadata, utils.ToBigEndianBytes(uint16(len(command)))...)
			metadata = append(metadata, command...)
		}
	}

	return metadata
}
