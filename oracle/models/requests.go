package models

import "github.com/holiman/uint256"

type OutboundTransactionItemRequest struct {
	Token           string `json:"token"`
	DistributionBPS int    `json:"distribution"`
	Recipient       string `json:"recipient"`
}

type InboundTransactionItemRequest struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type TransactionRequest struct {
	User        string                                           `json:"user"`
	SrcItems    map[Blockchain][]*InboundTransactionItemRequest  `json:"src_items"`
	DstItems    map[Blockchain][]*OutboundTransactionItemRequest `json:"dst_items"`
	SlippageBPS int                                              `json:"slippage_bps"`
}

type TransactionItemResponse struct {
	Token        string       `json:"token"`
	OutAmount    *uint256.Int `json:"out_amount"`
	MinOutAmount *uint256.Int `json:"min_out_amount"`
	USDXToken    *uint256.Int `json:"usdx_amount"`
}

type MarketplaceResponse struct {
	CommissionUSDX         uint64 `json:"commission_usdx"`
	CommissionDiscountUSDX uint64 `json:"commission_discount_usdx"`
	GasPriceUSDX           uint64 `json:"gas_price_usdx"`
	GasPriceDiscountUSDX   uint64 `json:"gas_price_discount_usdx"`
}

type UnsignedTransactionResponse struct {
	Contract string       `json:"contract"`
	Data     string       `json:"data"`
	Value    *uint256.Int `json:"value"`
}
type TransactionResponse struct {
	Marketplace     *MarketplaceResponse                      `json:"marketplace"`
	InTransaction   map[Blockchain][]*TransactionItemResponse `json:"in_transactions"`
	OutTransactions map[Blockchain][]*TransactionItemResponse `json:"out_transactions"`
	Transaction     *UnsignedTransactionResponse              `json:"transaction"`
}
