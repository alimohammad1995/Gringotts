package connection

type EstimateInboundTransfer struct {
	AmountUsdx uint64 `borsh:"amount_usdx"`
}

type EstimateOutboundTransferItem struct {
	Asset                   [32]byte `borsh:"asset"`
	ExecutionGasAmount      uint64   `borsh:"execution_gas"`
	ExecutionCommandLength  uint16   `borsh:"command_length"`
	ExecutionMetadataLength uint16   `borsh:"metadata_length"`
}

type EstimateOutboundTransfer struct {
	ChainId        uint8                          `borsh:"chain_id"`
	MetadataLength uint16                         `borsh:"metadata_length"`
	Items          []EstimateOutboundTransferItem `borsh:"items"`
}

type EstimateRequest struct {
	Inbound   EstimateInboundTransfer    `borsh:"inbound"`
	Outbounds []EstimateOutboundTransfer `borsh:"outbounds"`
}

type EstimateOutboundMetadata struct {
	ChainID                uint8  `borsh:"chain_id"`
	ExecutionGasAmount     uint64 `borsh:"execution_gas"`
	ExecutionGasAmountUsdx uint64 `borsh:"execution_gas_usdx"`
	TransferGasAmount      uint64 `borsh:"transfer_gas"`
	TransferGasAmountUsdx  uint64 `borsh:"transfer_gas_usdx"`
}

type EstimateResponse struct {
	CommissionUsdx          uint64                     `borsh:"commission_usdx"`
	CommissionDiscountUsdx  uint64                     `borsh:"commission_discount_usdx"`
	TransferGasPriceUsdx    uint64                     `borsh:"transfer_gas_usdx"`
	TransferGasDiscountUsdx uint64                     `borsh:"transfer_gas_discount_usdx"`
	OutboundMetadata        []EstimateOutboundMetadata `borsh:"outbound_metadata"`
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

type SolanaTransferItem struct {
	SwapAccountCount uint8  `borsh:"swap_account_count"`
	SwapCommand      []byte `borsh:"swap_command"`
}

type SolanaTransfer struct {
	AccountsAddress [][32]byte            `borsh:"accounts_address"`
	AccountFlags    []byte                `borsh:"accounts_flags"`
	AccountMapping  []byte                `borsh:"accounts_mapping"`
	Items           []*SolanaTransferItem `borsh:"items"`
}
