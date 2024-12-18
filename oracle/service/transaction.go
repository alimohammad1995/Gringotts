package service

import (
	"errors"
	"github.com/holiman/uint256"
	"gringotts/config"
	"gringotts/models"
	"gringotts/provider"
	"gringotts/utils"
)

func InboundTransaction(
	chain models.Blockchain,
	tokenAddress string,
	amount *uint256.Int,
	slippageBPS int,
) (*provider.Transaction, error) {
	token := models.GetToken(chain, tokenAddress)

	if token.IsStableCoin {
		return &provider.Transaction{
			FromToken:    token,
			ToToken:      token,
			OutAmount:    amount,
			MinOutAmount: amount,
		}, nil
	}

	var dex provider.Provider

	switch chain {
	case models.TRON:
		dex = &provider.SunSwap{}
	case models.Solana:
		dex = &provider.Jupiter{}
	default:
		dex = &provider.OpenOcean{}
	}

	stableCoin := models.GetDefaultStableCoins(chain)

	params := &provider.SwapParams{
		Chain:       chain,
		Amount:      amount,
		Recipient:   chain.GetContract(),
		SlippageBPS: slippageBPS,
		FromToken:   token,
		ToToken:     stableCoin,
	}

	swap, err := dex.GetSwap(params)
	if err != nil {
		return nil, err
	}
	if len(swap.Command) > config.MaxCommandLength {
		return nil, errors.New("swap command too long")
	}

	return &provider.Transaction{
		FromToken:    token,
		ToToken:      stableCoin,
		OutAmount:    swap.OutAmount,
		MinOutAmount: swap.MinOutAmount,
		Swap:         swap,
	}, nil
}

func OutboundTransaction(
	chain models.Blockchain,
	recipient string,
	tokenAddress string,
	amount *uint256.Int,
	slippageBPS int,
) (*provider.Transaction, error) {
	desiredToken := models.GetToken(chain, tokenAddress)

	if desiredToken.IsStableCoin {
		return &provider.Transaction{
			FromToken:    desiredToken,
			ToToken:      desiredToken,
			OutAmount:    utils.MoveDecimals(amount, config.ChainTransferDecimals, desiredToken.Decimals),
			MinOutAmount: utils.MoveDecimals(amount, config.ChainTransferDecimals, desiredToken.Decimals),
		}, nil
	}

	var dex provider.Provider
	switch chain {
	case models.TRON:
		dex = &provider.SunSwap{}
	case models.Solana:
		dex = &provider.Jupiter{}
	default:
		dex = &provider.OpenOcean{}
	}

	stableCoin := models.GetDefaultStableCoins(chain)
	stableCoinAmount := utils.MoveDecimals(amount, config.ChainTransferDecimals, stableCoin.Decimals)

	params := &provider.SwapParams{
		Chain:       chain,
		Amount:      stableCoinAmount,
		Recipient:   recipient,
		SlippageBPS: slippageBPS,
		FromToken:   stableCoin,
		ToToken:     desiredToken,
	}

	swap, err := dex.GetSwap(params)
	if err != nil {
		return nil, err
	}
	if len(swap.Command) > config.MaxCommandLength {
		return nil, errors.New("dex command too long")
	}

	return &provider.Transaction{
		FromToken:    stableCoin,
		ToToken:      desiredToken,
		Recipient:    recipient,
		OutAmount:    swap.OutAmount,
		MinOutAmount: swap.MinOutAmount,
		Swap:         swap,
	}, nil
}
