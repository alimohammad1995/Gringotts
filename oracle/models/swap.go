package models

import (
	"github.com/holiman/uint256"
)

type Side int

type Swap struct {
	Address       string   `json:"address"`
	Command       string   `json:"command"`
	MetaData      string   `json:"metadata"`
	AddressLookup []string `json:"addressLookup"`
}

type Transaction struct {
	FromToken string `json:"from_token"`
	ToToken   string `json:"to_token"`

	SrcAmount string

	Recipient       string `json:"recipient"`
	DistributionBPS BPS    `json:"distribution_bps"`

	OutAmount    *uint256.Int `json:"out_amount"`
	MinOutAmount *uint256.Int `json:"min_out_amount"`

	Swap *Swap `json:"swap"`
}

type Marketplace struct {
	Commission  float64 `json:"commission_usd"`
	GasPriceUSD float64 `json:"gas_price_usd"`
}

type UnsignedTransaction struct {
	Contract string `json:"contract"`
	Data     []byte `json:"data"`
	Value    string `json:"value"`
}
