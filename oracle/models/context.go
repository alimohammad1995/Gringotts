package models

import (
	"errors"
	"github.com/holiman/uint256"
	"gringotts/config"
	"gringotts/utils"
)

type Inbound struct {
	Amount      *uint256.Int
	Token       *Token
	SlippageBPS int

	Swap *Swap
}

type Outbound struct {
	Token           *Token
	DistributionBPS int
	Recipient       string
	SlippageBPS     int
	Amount          *uint256.Int
	Swap            *Swap
}

type ServingContext struct {
	Wallet string

	SourceChain Blockchain
	Inbounds    []*Inbound
	Outbounds   map[Blockchain][]*Outbound

	Marketplace *Marketplace
	Tx          *UnsignedTransaction
}

func NewServingContext(request *TransactionRequest) (*ServingContext, error) {
	servingContext := &ServingContext{
		Wallet: request.User,
		Marketplace: &Marketplace{
			GasPriceUSDX:           uint256.NewInt(0),
			GasPriceDiscountUSDX:   uint256.NewInt(0),
			CommissionUSDX:         uint256.NewInt(0),
			CommissionDiscountUSDX: uint256.NewInt(0),
		},
		Tx: &UnsignedTransaction{},
	}

	slippage := utils.Max(request.SlippageBPS, config.DefaultSlippageBPS)

	for chain, items := range request.SrcItems {
		if !chain.IsSupported() {
			return nil, errors.New("chain is not supported")
		}

		inbounds := make([]*Inbound, len(items))
		for i, item := range items {
			amount, err := uint256.FromDecimal(item.Amount)
			if err != nil {
				return nil, err
			}

			token := GetToken(chain, item.Token)
			if token == nil {
				return nil, errors.New("token does not exist")
			}

			inbounds[i] = &Inbound{
				Amount:      amount,
				Token:       token,
				SlippageBPS: slippage,
			}
		}

		servingContext.SourceChain = chain
		servingContext.Inbounds = inbounds
	}

	blockchainOutbounds := make(map[Blockchain][]*Outbound)
	for chain, items := range request.DstItems {
		if !chain.IsSupported() {
			return nil, errors.New("chain is not supported")
		}

		blockchainOutbounds[chain] = make([]*Outbound, len(items))
		for i, item := range items {
			token := GetToken(chain, item.Token)
			if token == nil {
				return nil, errors.New("token does not exist")
			}

			blockchainOutbounds[chain][i] = &Outbound{
				Token:           token,
				DistributionBPS: item.DistributionBPS,
				Recipient:       item.Recipient,
				SlippageBPS:     slippage,
			}
		}
	}
	servingContext.Outbounds = blockchainOutbounds

	if len(blockchainOutbounds) == 0 {
		return nil, errors.New("no blockchain outbound found")
	}

	return servingContext, nil
}
