package service

import (
	"context"
	solanaClient "github.com/blocto/solana-go-sdk/client"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/address_lookup_table"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gofiber/fiber/v2/log"
	"github.com/holiman/uint256"
	"github.com/near/borsh-go"
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
	blockchain models.Blockchain,
	inTransactions map[models.Blockchain][]*models.Transaction,
	outTransactions map[models.Blockchain][]*models.Transaction,
) (*provider.UnsignedTransaction, error) {
	inTransaction := inTransactions[blockchain]
	alts := make([]string, 0)

	inboundTransfers := make([]connection.BridgeInboundTransferItem, len(inTransaction))
	for i, transaction := range inTransaction {
		tAmount, _ := uint256.FromDecimal(transaction.SrcAmount)
		inboundTransfer := connection.BridgeInboundTransferItem{
			Asset:  utils.ToByte32(transaction.FromToken),
			Amount: tAmount.Uint64(),
		}

		if transaction.Swap != nil {
			inboundTransfer.Swap = &connection.Swap{
				Executor:    utils.ToByte32(transaction.Swap.Address),
				Command:     utils.FromHex(transaction.Swap.Command),
				Metadata:    utils.FromHex(transaction.Swap.MetaData),
				StableToken: utils.ToByte32(transaction.ToToken),
			}

			alts = append(alts, transaction.Swap.AddressLookup...)
		}

		inboundTransfers[i] = inboundTransfer
	}

	outboundTransfers := make([]connection.BridgeOutboundTransfer, 0, len(outTransactions))
	for chain, transactions := range outTransactions {
		outboundTransferItems := make([]connection.BridgeOutboundTransferItem, len(transactions))

		for i, transaction := range transactions {
			gas, _, _ := GetExecutionParams(chain, transaction.ToToken)

			outboundTransferItem := connection.BridgeOutboundTransferItem{
				Asset:              utils.ToByte32(transaction.ToToken),
				Recipient:          utils.ToByte32(transaction.Recipient),
				ExecutionGasAmount: gas,
				DistributionBP:     uint16(transaction.DistributionBPS),
			}

			if transaction.Swap != nil {
				outboundTransferItem.Swap = &connection.Swap{
					Executor:    utils.ToByte32(transaction.Swap.Address),
					Command:     utils.FromHex(transaction.Swap.Command),
					Metadata:    utils.FromHex(transaction.Swap.MetaData),
					StableToken: utils.ToByte32(transaction.FromToken),
				}
			}

			outboundTransferItems[i] = outboundTransferItem
		}

		outboundTransfers = append(outboundTransfers, connection.BridgeOutboundTransfer{
			ChainID: chain.GetId(),
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

	accounts := []types.AccountMeta{
		{PubKey: solana.PublicKeyFromString(models.GetPriceFeed()), IsSigner: false, IsWritable: false},
		{PubKey: solana.PublicKeyFromString(models.GetGringotts(blockchain)), IsSigner: false, IsWritable: false},
	}

	instruction := types.Instruction{
		ProgramID: solana.PublicKeyFromString(blockchain.GetContract()),
		Accounts:  accounts,
		Data:      data,
	}

	client := solanaClient.NewClient(blockchain.GetEndpoint())
	recentBlockhash, _ := client.GetLatestBlockhash(context.Background())

	solanaALTs := make([]types.AddressLookupTableAccount, 0)
	for _, alt := range alts {
		accountInfo, err := client.GetAccountInfo(context.Background(), alt)
		if err != nil {
			return nil, err
		}
		addressLookupTable, err := address_lookup_table.DeserializeLookupTable(accountInfo.Data, accountInfo.Owner)
		if err != nil {
			log.Errorw("invalid alt", "err", err)
			continue
		}
		solanaALTs = append(solanaALTs, types.AddressLookupTableAccount{
			Key:       solana.PublicKeyFromString(alt),
			Addresses: addressLookupTable.Addresses,
		})
	}

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:                   solana.PublicKeyFromString(wallet),
		RecentBlockhash:            recentBlockhash.Blockhash,
		Instructions:               []types.Instruction{instruction},
		AddressLookupTableAccounts: solanaALTs,
	})

	tx, _ := message.Serialize()

	return &provider.UnsignedTransaction{
		Contract: blockchain.GetContract(),
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
				Executor:    utils.ToByte32(transaction.Swap.Address),
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
				Asset:              utils.ToByte32(transaction.ToToken),
				Recipient:          utils.ToByte32(transaction.Recipient),
				ExecutionGasAmount: big.NewInt(int64(gas)),
				DistributionBP:     uint16(transaction.DistributionBPS),
			}

			if transaction.Swap != nil {
				transactionItem.Swap = connection.GringottsSwap{
					Executor:    utils.ToByte32(transaction.Swap.Address),
					Command:     utils.FromHex(transaction.Swap.Command),
					Metadata:    utils.FromHex(transaction.Swap.MetaData),
					StableToken: utils.ToByte32(transaction.FromToken),
				}
			}

			transactionItems[i] = transactionItem
		}

		outboundTransfers = append(outboundTransfers, connection.GringottsBridgeOutboundTransfer{
			ChainId: chain.GetId(),
			Items:   transactionItems,
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
