package service

import (
	"errors"
	"github.com/holiman/uint256"
	"gringotts/config"
	"gringotts/models"
	"gringotts/provider"
	"gringotts/utils"
)

func InboundTransaction(chain models.Blockchain, item *models.Inbound) error {
	swap, err := InboundTransactionSwap(chain, item.Token, item.Amount, item.SlippageBPS)

	if err != nil {
		return err
	}

	item.Swap = swap
	return nil
}

func InboundTransactionSwap(chain models.Blockchain, token *models.Token, amount *uint256.Int, slippageBPS int) (*models.Swap, error) {
	if token.IsStableCoin {
		return nil, nil
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

	params := &models.SwapParams{
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

	swap.FromToken = token
	swap.ToToken = stableCoin

	return swap, nil
}

func OutboundTransaction(chain models.Blockchain, item *models.Outbound) error {
	swap, err := OutboundTransactionSwap(chain, item.Recipient, item.Token, item.Amount, item.SlippageBPS)

	if err != nil {
		return err
	}

	item.Swap = swap
	return nil
}

func OutboundTransactionSwap(
	chain models.Blockchain,
	recipient string,
	token *models.Token,
	amount *uint256.Int,
	slippageBPS int,
) (*models.Swap, error) {
	if token.IsStableCoin {
		return nil, nil
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

	params := &models.SwapParams{
		Chain:       chain,
		Amount:      stableCoinAmount,
		Recipient:   recipient,
		SlippageBPS: slippageBPS,
		FromToken:   stableCoin,
		ToToken:     token,
	}

	if chain.IsSolana() && token.IsNative {
		params.Recipient = models.GetGringotts(chain)
	}

	swap, err := dex.GetSwap(params)
	if err != nil {
		return nil, err
	}
	if len(swap.Command) > config.MaxCommandLength {
		return nil, errors.New("dex command too long")
	}

	swap.FromToken = stableCoin
	swap.ToToken = token

	return swap, nil
}
