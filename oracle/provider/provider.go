package provider

import (
	"github.com/holiman/uint256"
	"gringotts/models"
)

type SwapParams struct {
	Chain       models.Blockchain
	FromToken   *models.Token
	ToToken     *models.Token
	Amount      *uint256.Int
	Recipient   string
	SlippageBPS int
}

type Swap struct {
	ExecutorAddress string
	Command         string
	Metadata        string
	AddressLookup   []string

	OutAmount    *uint256.Int
	MinOutAmount *uint256.Int
}

type Transaction struct {
	FromToken *models.Token
	ToToken   *models.Token

	SrcAmount *uint256.Int

	Recipient       string
	DistributionBPS int

	OutAmount    *uint256.Int
	MinOutAmount *uint256.Int

	Swap *Swap
}

type Estimate struct {
	GasPrice      *uint256.Int
	GasPriceUSD   *uint256.Int
	CommissionUSD *uint256.Int
}

type UnsignedTransaction struct {
	Contract string
	Data     []byte
	Value    string
}

type Provider interface {
	GetSwap(params *SwapParams) (*Swap, error)
}
