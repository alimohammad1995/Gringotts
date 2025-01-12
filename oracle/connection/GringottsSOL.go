package connection

type EstimateInboundTransfer struct {
	AmountUsdx uint64 `borsh:"amount_usdx"`
}
type EstimateOutboundTransfer struct {
	ChainId       uint8  `borsh:"chain_id"`
	ExecutionGas  uint64 `borsh:"execution_gas"`
	MessageLength uint16 `borsh:"message_length"`
}

type EstimateRequest struct {
	Inbound   EstimateInboundTransfer    `borsh:"inbound"`
	Outbounds []EstimateOutboundTransfer `borsh:"outbounds"`
}

type EstimateOutboundDetails struct {
	ChainID            uint8  `borsh:"chain_id"`
	ExecutionGasAmount uint64 `borsh:"execution_gas"`
	TransferGasAmount  uint64 `borsh:"transfer_gas"`
}

type EstimateResponse struct {
	CommissionUsdx          uint64                    `borsh:"commission_usdx"`
	CommissionDiscountUsdx  uint64                    `borsh:"commission_discount_usdx"`
	TransferGasPriceUsdx    uint64                    `borsh:"transfer_gas_usdx"`
	TransferGasDiscountUsdx uint64                    `borsh:"transfer_gas_discount_usdx"`
	OutboundDetails         []EstimateOutboundDetails `borsh:"outbound_details"`
}

type Swap struct {
	Executor      [32]byte `borsh:"executor"`
	Command       []byte   `borsh:"command"`
	AccountsCount uint8    `borsh:"accounts_count"`
	StableToken   [32]byte `borsh:"stable_token"`
}

type BridgeInboundTransferItem struct {
	Asset  [32]byte `borsh:"asset"`
	Amount uint64   `borsh:"amount"`
	Swap   *Swap    `borsh:"swap"`
}

type BridgeInboundTransfer struct {
	AmountUSDx uint64                      `borsh:"amount_usdx"`
	Items      []BridgeInboundTransferItem `borsh:"items"`
}

type BridgeOutboundTransferItem struct {
	DistributionBP uint16 `borsh:"distribution_bp"`
}

type BridgeOutboundTransfer struct {
	ChainID      uint8                        `borsh:"chain_id"`
	ExecutionGas uint64                       `borsh:"execution_gas"`
	Message      []byte                       `borsh:"message"`
	Items        []BridgeOutboundTransferItem `borsh:"items"`
}

type BridgeRequest struct {
	Inbound   BridgeInboundTransfer    `borsh:"inbound"`
	Outbounds []BridgeOutboundTransfer `borsh:"outbounds"`
}
