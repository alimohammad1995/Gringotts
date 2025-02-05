package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	solana_client "github.com/blocto/solana-go-sdk/client"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/near/borsh-go"
	"gringotts/blockchain"
	"gringotts/config"
	"gringotts/connection"
	"gringotts/models"
)

func EstimateMarketplace(
	servingContext *models.ServingContext,
	amount *uint256.Int,
) error {
	var res *models.Marketplace
	var err error

	switch servingContext.SourceChain {
	case models.Solana, models.SolanaDev:
		res, err = EstimateMarketplaceSolana(servingContext, amount)
	default:
		res, err = EstimateMarketplaceEVM(servingContext, amount)
	}

	if err != nil {
		return err
	}

	servingContext.Marketplace = res
	return nil
}

func EstimateMarketplaceSolana(
	servingContext *models.ServingContext,
	amount *uint256.Int,
) (*models.Marketplace, error) {
	inbound := connection.EstimateInboundTransfer{
		AmountUsdx: amount.Uint64(),
	}
	outbounds := make([]connection.EstimateOutboundTransfer, 0, len(servingContext.Outbounds))
	for chain, items := range servingContext.Outbounds {
		gas := uint64(0)
		totalMessageLength := uint16(0)

		for _, item := range items {
			executionGas, commendLength := GetExecutionParams(chain, item.Token)

			gas = gas + executionGas
			totalMessageLength = totalMessageLength + commendLength
		}

		outbounds = append(outbounds, connection.EstimateOutboundTransfer{
			ChainId:       chain.GetId(),
			ExecutionGas:  gas,
			MessageLength: totalMessageLength,
		})
	}

	requestSerializedData, _ := borsh.Serialize(connection.EstimateRequest{
		Inbound:   inbound,
		Outbounds: outbounds,
	})
	data := append(models.GetEstimateDiscriminator(), requestSerializedData...)

	accounts := []types.AccountMeta{
		{PubKey: solana.PublicKeyFromString(models.GetPriceFeed()), IsSigner: false, IsWritable: false},
		{PubKey: solana.PublicKeyFromString(models.GetGringotts(servingContext.SourceChain)), IsSigner: false, IsWritable: false},
	}
	for chain := range servingContext.Outbounds {
		accounts = append(accounts, types.AccountMeta{
			PubKey: solana.PublicKeyFromString(models.GetPeer(servingContext.SourceChain, chain)), IsSigner: false, IsWritable: false,
		})
	}
	for chain := range servingContext.Outbounds {
		accounts = append(accounts, getEstimateAccounts(chain)...)
	}

	instruction := types.Instruction{
		ProgramID: solana.PublicKeyFromString(servingContext.SourceChain.GetContract()),
		Accounts:  accounts,
		Data:      data,
	}

	client := solana_client.NewClient(servingContext.SourceChain.GetEndpoint())
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

	programResult := fmt.Sprintf("Program return: %s ", servingContext.SourceChain.GetContract())
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

	return &models.Marketplace{
		GasPriceUSDX:           uint256.NewInt(response.TransferGasPriceUsdx),
		GasPriceDiscountUSDX:   uint256.NewInt(response.TransferGasDiscountUsdx),
		CommissionUSDX:         uint256.NewInt(response.CommissionUsdx),
		CommissionDiscountUSDX: uint256.NewInt(response.CommissionDiscountUsdx),
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

func EstimateMarketplaceEVM(
	servingContext *models.ServingContext,
	amount *uint256.Int,
) (*models.Marketplace, error) {
	client, _ := blockchain.GetConnection(servingContext.SourceChain)
	instance, _ := connection.NewGringottsEVMCaller(common.HexToAddress(servingContext.SourceChain.GetContract()), client)

	inbound := connection.GringottsEstimateInboundTransfer{
		AmountUSDX: amount.ToBig(),
	}

	outbounds := make([]connection.GringottsEstimateOutboundTransfer, 0, len(servingContext.Outbounds))
	for chain, items := range servingContext.Outbounds {
		gas := int64(0)
		totalMessageLength := uint16(0)

		for _, item := range items {
			executionGas, commendLength := GetExecutionParams(chain, item.Token)

			gas = gas + int64(executionGas)
			totalMessageLength = totalMessageLength + commendLength
		}

		outbounds = append(outbounds, connection.GringottsEstimateOutboundTransfer{
			ChainId:       chain.GetId(),
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
	gasPriceDiscountUSD, _ := uint256.FromBig(res.TransferGasDiscountUSDX)
	commissionUSD, _ := uint256.FromBig(res.CommissionUSDX)
	commissionDiscountUSD, _ := uint256.FromBig(res.CommissionDiscountUSDX)

	return &models.Marketplace{
		GasPriceUSDX:           gasPriceUSD,
		GasPriceDiscountUSDX:   gasPriceDiscountUSD,
		CommissionUSDX:         commissionUSD,
		CommissionDiscountUSDX: commissionDiscountUSD,
	}, nil
}

// GetExecutionParams returns gas, message length
func GetExecutionParams(chain models.Blockchain, token *models.Token) (uint64, uint16) {
	if token.IsStableCoin {
		switch chain {
		case models.Solana, models.SolanaDev:
			return 100_000, 0
		default:
			return 100_000, 0
		}
	}

	switch chain {
	case models.Solana, models.SolanaDev:
		return 500_000, config.MaxCommandLength
	default:
		return 500_000, config.MaxCommandLength
	}
}
