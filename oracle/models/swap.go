package models

import (
	"github.com/holiman/uint256"
)

type Swap struct {
	Executor      string    `json:"executor"`
	Command       string    `json:"command"`
	Accounts      []Account `json:"accounts"`
	AddressLookup []string  `json:"addressLookup"`
}

type Transaction struct {
	FromToken string `json:"from_token"`
	ToToken   string `json:"to_token"`

	SrcAmount string `json:"src_amount"`

	Recipient       string `json:"recipient"`
	DistributionBPS int    `json:"distribution_bps"`

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
