package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"gringotts/blockchain"
	"gringotts/config"
	"gringotts/connection"
	"gringotts/models"
	"gringotts/provider"
	"gringotts/utils"

	solana_client "github.com/blocto/solana-go-sdk/client"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/near/borsh-go"
)

func EstimateMarketplace(
	chain models.Blockchain,
	dstItems map[models.Blockchain][]string,
	amount *uint256.Int,
) (*provider.Estimate, error) {
	switch chain {
	case models.Solana, models.SolanaDev:
		return EstimateMarketplaceSolana(chain, dstItems, amount)
	default:
		return estimateMarketplaceEVM(chain, dstItems, amount)
	}
}

func EstimateMarketplaceSolana(
	chain models.Blockchain,
	dstItems map[models.Blockchain][]string,
	amount *uint256.Int,
) (*provider.Estimate, error) {
	inbound := connection.EstimateInboundTransfer{
		AmountUsdx: amount.Uint64(),
	}
	outbounds := make([]connection.EstimateOutboundTransfer, 0, len(dstItems))
	for chainIter, assets := range dstItems {
		items := make([]connection.EstimateOutboundTransferItem, 0, len(assets))

		for _, asset := range assets {
			executionGas, commendLength, metadataLength := GetExecutionParams(chainIter, asset)

			items = append(items, connection.EstimateOutboundTransferItem{
				Asset:                   utils.ToByte32(asset),
				ExecutionGasAmount:      executionGas,
				ExecutionCommandLength:  commendLength,
				ExecutionMetadataLength: metadataLength,
			})
		}

		outbounds = append(outbounds, connection.EstimateOutboundTransfer{
			ChainId: chainIter.GetId(),
			Items:   items,
		})
	}

	requestSerializedData, _ := borsh.Serialize(connection.EstimateRequest{
		Inbound:   inbound,
		Outbounds: outbounds,
	})
	data := append(models.GetEstimateDiscriminator(), requestSerializedData...)

	accounts := []types.AccountMeta{
		{PubKey: solana.PublicKeyFromString(models.GetPriceFeed()), IsSigner: false, IsWritable: false},
		{PubKey: solana.PublicKeyFromString(models.GetGringotts(chain)), IsSigner: false, IsWritable: false},
	}
	for chainIter := range dstItems {
		accounts = append(accounts, types.AccountMeta{
			PubKey: solana.PublicKeyFromString(models.GetPeer(chain, chainIter)), IsSigner: false, IsWritable: false,
		})
	}
	for chainIter := range dstItems {
		accounts = append(accounts, getEstimateAccounts(chainIter)...)
	}

	instruction := types.Instruction{
		ProgramID: solana.PublicKeyFromString(chain.GetContract()),
		Accounts:  accounts,
		Data:      data,
	}

	client := solana_client.NewClient(chain.GetEndpoint())
	recentBlockhash, _ := client.GetLatestBlockhash(context.Background())

	signer, _ := types.AccountFromBase58(models.GetSigner())

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        signer.PublicKey,
			RecentBlockhash: recentBlockhash.Blockhash,
			Instructions:    []types.Instruction{instruction},
		}),
		Signers: []types.Account{signer},
	})
	if err != nil {
		return nil, err
	}

	simResult, err := client.SimulateTransaction(context.Background(), tx)
	if err != nil {
		return nil, err
	}

	programResult := fmt.Sprintf("Program return: %s ", chain.GetContract())
	var output string
	for _, log := range simResult.Logs {
		if strings.HasPrefix(log, programResult) {
			output = log[len(programResult):]
		}
	}

	decodedData, err := base64.StdEncoding.DecodeString(output)
	var response connection.EstimateResponse
	err = borsh.Deserialize(&response, decodedData)
	if err != nil {
		return nil, err
	}

	return &provider.Estimate{
		GasPriceUSD:   uint256.NewInt(response.TransferGasPriceUsdx),
		CommissionUSD: uint256.NewInt(response.CommissionUsdx),
	}, nil
}

func getEstimateAccounts(chain models.Blockchain) []types.AccountMeta {
	accounts := models.GetEstimateAccounts(chain)
	res := make([]types.AccountMeta, len(accounts))

	for i, account := range accounts {
		res[i] = types.AccountMeta{
			PubKey:     solana.PublicKeyFromString(account.Address),
			IsSigner:   account.IsSigner,
			IsWritable: account.IsWritable,
		}
	}

	return res
}

func estimateMarketplaceEVM(
	chain models.Blockchain,
	dstItems map[models.Blockchain][]string,
	amount *uint256.Int,
) (*provider.Estimate, error) {
	client, _ := blockchain.GetConnection(chain)
	instance, _ := connection.NewGringottsEVMCaller(common.HexToAddress(chain.GetContract()), client)

	inbound := connection.GringottsEstimateInboundTransfer{
		AmountUSDX: amount.ToBig(),
	}

	outbounds := make([]connection.GringottsEstimateOutboundTransfer, 0, len(dstItems))
	for chainIter, assets := range dstItems {
		gas := int64(0)
		totalMessageLength := uint16(0)

		for _, asset := range assets {
			executionGas, commendLength, metadataLength := GetExecutionParams(chainIter, asset)

			gas = gas + int64(executionGas)
			totalMessageLength = totalMessageLength + commendLength + metadataLength
		}

		outbounds = append(outbounds, connection.GringottsEstimateOutboundTransfer{
			ChainId:       chainIter.GetId(),
			MessageLength: totalMessageLength,
			ExecutionGas:  big.NewInt(gas),
		})
	}

	res, err := instance.Estimate(&bind.CallOpts{}, connection.GringottsEstimateRequest{
		Inbound:   inbound,
		Outbounds: outbounds,
	})
	if err != nil {
		return nil, err
	}

	gasPriceUSD, _ := uint256.FromBig(res.TransferGasUSDX)
	CommissionUSD, _ := uint256.FromBig(res.CommissionUSDX)

	return &provider.Estimate{
		GasPriceUSD:   gasPriceUSD,
		CommissionUSD: CommissionUSD,
	}, nil
}

func GetExecutionParams(chain models.Blockchain, asset string) (uint64, uint16, uint16) {
	token := models.GetToken(chain, asset)

	if token.IsStableCoin {
		switch chain {
		case models.Solana, models.SolanaDev:
			return 100_000, 0, 0
		default:
			return 100_000, 0, 0
		}
	}

	switch chain {
	case models.Solana, models.SolanaDev:
		return 500_000, config.MaxCommandLength, config.MaxMetaDataLength
	default:
		return 500_000, config.MaxCommandLength, 0
	}
}
