package service

import (
	"context"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/btcsuite/btcutil/base58"
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
			Asset:  utils.ToByte32(transaction.FromToken),
			Amount: tAmount.Uint64(),
		}

		if transaction.Swap != nil {
			inboundTransfer.Swap = &connection.Swap{
				Executor:    utils.ToByte32(transaction.Swap.Executor),
				Command:     utils.FromHex(transaction.Swap.Command),
				Metadata:    utils.FromHex(transaction.Swap.MetaData),
				StableToken: utils.ToByte32(transaction.ToToken),
			}
		}

		inboundTransfers[i] = inboundTransfer
	}

	outboundTransfers := make([]connection.BridgeOutboundTransfer, 0, len(outTransactions))
	for chainIter, transactions := range outTransactions {
		outboundTransferItems := make([]connection.BridgeOutboundTransferItem, len(transactions))

		for i, transaction := range transactions {
			gas, _, _ := GetExecutionParams(chainIter, transaction.ToToken)

			outboundTransferItem := connection.BridgeOutboundTransferItem{
				Asset:              utils.ToByte32(transaction.ToToken),
				Recipient:          utils.ToByte32(transaction.Recipient),
				ExecutionGasAmount: gas,
				DistributionBP:     uint16(transaction.DistributionBPS),
			}

			if transaction.Swap != nil {
				outboundTransferItem.Swap = &connection.Swap{
					Executor:    utils.ToByte32(transaction.Swap.Executor),
					Command:     utils.FromHex(transaction.Swap.Command),
					Metadata:    utils.FromHex(transaction.Swap.MetaData),
					StableToken: utils.ToByte32(transaction.FromToken),
				}
			}

			outboundTransferItems[i] = outboundTransferItem
		}

		outboundTransfers = append(outboundTransfers, connection.BridgeOutboundTransfer{
			ChainID: chainIter.GetId(),
			Items:   outboundTransferItems,
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
			transaction.Swap.MetaData = utils.ToHex([]byte{byte(len(transaction.Swap.Accounts))})
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
			transactionItem.Swap = connection.GringottsSwap{
				Executor:    utils.ToByte32(transaction.Swap.Executor),
				Command:     utils.FromHex(transaction.Swap.Command),
				Metadata:    utils.FromHex(transaction.Swap.MetaData),
				StableToken: utils.ToByte32(transaction.ToToken),
			}
		}

		if len(transaction.FromToken) == 0 {
			value = transaction.SrcAmount
		}

		inboundTransfers[i] = transactionItem
	}

	outboundTransfers := make([]connection.GringottsBridgeOutboundTransfer, 0, len(outTransactions))
	for chain, transactions := range outTransactions {
		transactionItems := make([]connection.GringottsBridgeOutboundTransferItem, len(transactions))

		for i, transaction := range transactions {
			gas, _, _ := GetExecutionParams(chain, transaction.ToToken)

			transactionItem := connection.GringottsBridgeOutboundTransferItem{
				Asset:          utils.ToByte32(transaction.ToToken),
				Recipient:      utils.ToByte32(transaction.Recipient),
				ExecutionGas:   big.NewInt(int64(gas)),
				DistributionBP: uint16(transaction.DistributionBPS),
			}

			if transaction.Swap != nil {
				transactionItem.Swap = connection.GringottsSwap{
					Executor:    utils.ToByte32(transaction.Swap.Executor),
					Command:     utils.FromHex(transaction.Swap.Command),
					Metadata:    utils.FromHex(transaction.Swap.MetaData),
					StableToken: utils.ToByte32(transaction.FromToken),
				}
			}

			transactionItems[i] = transactionItem
		}

		outboundTransfers = append(outboundTransfers, connection.GringottsBridgeOutboundTransfer{
			ChainId:  chain.GetId(),
			Metadata: getMetaData(chain, transactions),
			Items:    transactionItems,
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

	return &provider.UnsignedTransaction{
		Contract: blockchain.GetContract(),
		Data:     data,
		Value:    value,
	}, nil
}

func getMetaData(chain models.Blockchain, transactions []*models.Transaction) []byte {
	if chain != models.Solana && chain != models.SolanaDev {
		return make([]byte, 0)
	}

	accounts := []models.Account{
		{Address: models.GetGringotts(chain), IsSigner: false, IsWritable: false},
		{Address: models.GetVault(chain), IsSigner: false, IsWritable: true},
		{Address: solana.SPLAssociatedTokenAccountProgramID.String(), IsSigner: false, IsWritable: false},
		{Address: solana.TokenProgramID.String(), IsSigner: false, IsWritable: false},
		{Address: solana.SystemProgramID.String(), IsSigner: false, IsWritable: false},
	}

	for _, transaction := range transactions {
		accounts = append(accounts,
			models.Account{Address: transaction.Recipient, IsSigner: false, IsWritable: transaction.ToToken == ""},
		)

		// Stable transfer
		if transaction.Swap == nil {
			accounts = append(accounts,
				models.Account{Address: transaction.ToToken, IsSigner: false, IsWritable: false},
				models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), transaction.ToToken), IsSigner: false, IsWritable: true},
				models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.ToToken), IsSigner: false, IsWritable: true},
			)
		} else {
			if transaction.ToToken != "" {
				accounts = append(accounts,
					models.Account{Address: transaction.FromToken, IsSigner: false, IsWritable: false},
					models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), transaction.FromToken), IsSigner: false, IsWritable: true},
					models.Account{Address: transaction.ToToken, IsSigner: false, IsWritable: false},
					models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.ToToken), IsSigner: false, IsWritable: true},
					models.Account{Address: provider.JupiterAddress, IsSigner: false, IsWritable: false},
					models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.FromToken), IsSigner: false, IsWritable: true},
				)
			} else {
				accounts = append(accounts,
					models.Account{Address: transaction.FromToken, IsSigner: false, IsWritable: false},
					models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), transaction.FromToken), IsSigner: false, IsWritable: true},
					models.Account{Address: models.GetAssociatedTokenAddress(models.GetGringotts(chain), models.NativeMint), IsSigner: false, IsWritable: true},
					models.Account{Address: provider.JupiterAddress, IsSigner: false, IsWritable: false},
					models.Account{Address: models.GetAssociatedTokenAddress(transaction.Recipient, transaction.FromToken), IsSigner: false, IsWritable: true},
				)
			}

			for _, acc := range transaction.Swap.Accounts {
				accounts = append(accounts,
					models.Account{Address: acc.Address, IsSigner: false, IsWritable: acc.IsWritable},
				)
			}
		}
	}

	accountsMap := make(map[string]models.Account)
	for _, account := range accounts {
		if existingAccount, ok := accountsMap[account.Address]; ok {
			existingAccount.IsWritable = accountsMap[account.Address].IsWritable || account.IsWritable
			accountsMap[account.Address] = existingAccount
		} else {
			accountsMap[account.Address] = account
		}
	}

	metadata := make([]byte, len(accountsMap))
	addressMap := make(map[string]int)
	i := 0
	for _, account := range accountsMap {
		addressBytes := base58.Decode(account.Address)
		metadata = append(metadata, addressBytes...)
		if account.IsWritable {
			metadata = append(metadata, 1)
		} else {
			metadata = append(metadata, 0)
		}
		addressMap[account.Address] = i
		i += 1
	}

	for _, account := range accounts {
		metadata = append(metadata, byte(addressMap[account.Address]))
	}

	return metadata
}
