package models

import "github.com/holiman/uint256"

type SwapParams struct {
	Chain       Blockchain
	FromToken   *Token
	ToToken     *Token
	Amount      *uint256.Int
	Recipient   string
	SlippageBPS int
}

type Account struct {
	Address    string
	IsSigner   bool
	IsWritable bool
}

type Swap struct {
	Executor      string
	Command       string
	Accounts      []*Account
	AddressLookup []string

	OutAmount    *uint256.Int
	MinOutAmount *uint256.Int

	FromToken *Token
	ToToken   *Token
}

type UnsignedTransaction struct {
	Data  []byte
	Value *uint256.Int
}

type Marketplace struct {
	GasPriceUSDX           *uint256.Int
	GasPriceDiscountUSDX   *uint256.Int
	CommissionUSDX         *uint256.Int
	CommissionDiscountUSDX *uint256.Int
}
