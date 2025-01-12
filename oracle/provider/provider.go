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

type Account struct {
	Address     string
	IsSigner    bool
	IsWriteable bool
}

type Swap struct {
	ExecutorAddress string
	Command         string
	Metadata        string
	Accounts        []Account
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
	GasPriceUSDX           *uint256.Int
	GasPriceDiscountUSDX   *uint256.Int
	CommissionUSDX         *uint256.Int
	CommissionDiscountUSDX *uint256.Int
}

type UnsignedTransaction struct {
	Contract string
	Data     []byte
	Value    string
}

type Provider interface {
	GetSwap(params *SwapParams) (*Swap, error)
}
