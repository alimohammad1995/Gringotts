// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package connection

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Config is an auto generated low-level Go binding around an user-defined struct.
type Config struct {
	CommissionMicroBPS    uint32
	CommissionDiscountBPS uint32
	GasDiscountBPS        uint32
	StableCoins           [][32]byte
}

// GringottsBridgeInboundTransfer is an auto generated low-level Go binding around an user-defined struct.
type GringottsBridgeInboundTransfer struct {
	AmountUSDX *big.Int
	Items      []GringottsBridgeInboundTransferItem
}

// GringottsBridgeInboundTransferItem is an auto generated low-level Go binding around an user-defined struct.
type GringottsBridgeInboundTransferItem struct {
	Asset       [32]byte
	Amount      *big.Int
	Executor    [32]byte
	Command     []byte
	StableToken [32]byte
}

// GringottsBridgeOutboundTransfer is an auto generated low-level Go binding around an user-defined struct.
type GringottsBridgeOutboundTransfer struct {
	ChainId      uint8
	ExecutionGas *big.Int
	Message      []byte
	Items        []GringottsBridgeOutboundTransferItem
}

// GringottsBridgeOutboundTransferItem is an auto generated low-level Go binding around an user-defined struct.
type GringottsBridgeOutboundTransferItem struct {
	DistributionBP uint16
}

// GringottsBridgeRequest is an auto generated low-level Go binding around an user-defined struct.
type GringottsBridgeRequest struct {
	Inbound   GringottsBridgeInboundTransfer
	Outbounds []GringottsBridgeOutboundTransfer
}

// GringottsBridgeResponse is an auto generated low-level Go binding around an user-defined struct.
type GringottsBridgeResponse struct {
	MessageIds [][32]byte
}

// GringottsEstimateInboundTransfer is an auto generated low-level Go binding around an user-defined struct.
type GringottsEstimateInboundTransfer struct {
	AmountUSDX *big.Int
}

// GringottsEstimateOutboundDetails is an auto generated low-level Go binding around an user-defined struct.
type GringottsEstimateOutboundDetails struct {
	ChainId      uint8
	ExecutionGas *big.Int
	TransferGas  *big.Int
}

// GringottsEstimateOutboundTransfer is an auto generated low-level Go binding around an user-defined struct.
type GringottsEstimateOutboundTransfer struct {
	ChainId       uint8
	ExecutionGas  *big.Int
	MessageLength uint16
}

// GringottsEstimateRequest is an auto generated low-level Go binding around an user-defined struct.
type GringottsEstimateRequest struct {
	Inbound   GringottsEstimateInboundTransfer
	Outbounds []GringottsEstimateOutboundTransfer
}

// GringottsEstimateResponse is an auto generated low-level Go binding around an user-defined struct.
type GringottsEstimateResponse struct {
	CommissionUSDX          *big.Int
	CommissionDiscountUSDX  *big.Int
	TransferGasUSDX         *big.Int
	TransferGasDiscountUSDX *big.Int
	OutboundDetails         []GringottsEstimateOutboundDetails
}

// Origin is an auto generated low-level Go binding around an user-defined struct.
type Origin struct {
	SrcEid uint32
	Sender [32]byte
	Nonce  uint64
}

// Peer is an auto generated low-level Go binding around an user-defined struct.
type Peer struct {
	ChainID               uint8
	Endpoint              [32]byte
	LzEID                 uint32
	MultiSend             bool
	BaseGasEstimate       *big.Int
	RegisterGasEstimate   *big.Int
	CompletionGasEstimate *big.Int
}

// GringottsEVMMetaData contains all meta data concerning the GringottsEVM contract.
var GringottsEVMMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"ChainId\",\"name\":\"_chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_networkDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_lzEndpoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BlockedSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDelegate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidEndpointCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"optionType\",\"type\":\"uint16\"}],\"name\":\"InvalidOptionType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LzTokenUnavailable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"eid\",\"type\":\"uint32\"}],\"name\":\"NoPeer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"msgValue\",\"type\":\"uint256\"}],\"name\":\"NotEnoughNative\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"OnlyEndpoint\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"eid\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"}],\"name\":\"OnlyPeer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintDowncast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"blocked\",\"type\":\"bool\"}],\"name\":\"BlockEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"eid\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"peer\",\"type\":\"bytes32\"}],\"name\":\"PeerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"ChainId\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"messageId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountUSDX\",\"type\":\"uint256\"}],\"name\":\"ReceiveChainTransferEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"ChainId\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"messageId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountUSDX\",\"type\":\"uint256\"}],\"name\":\"SendChainTransferEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"}],\"name\":\"TestMessage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"srcEid\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structOrigin\",\"name\":\"origin\",\"type\":\"tuple\"}],\"name\":\"allowInitializePath\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"balanceERC20\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"blockAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amountUSDX\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"GringottsAddress\",\"name\":\"asset\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"GringottsAddress\",\"name\":\"executor\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"command\",\"type\":\"bytes\"},{\"internalType\":\"GringottsAddress\",\"name\":\"stableToken\",\"type\":\"bytes32\"}],\"internalType\":\"structGringotts.BridgeInboundTransferItem[]\",\"name\":\"items\",\"type\":\"tuple[]\"}],\"internalType\":\"structGringotts.BridgeInboundTransfer\",\"name\":\"inbound\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"ChainId\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"executionGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint16\",\"name\":\"distributionBP\",\"type\":\"uint16\"}],\"internalType\":\"structGringotts.BridgeOutboundTransferItem[]\",\"name\":\"items\",\"type\":\"tuple[]\"}],\"internalType\":\"structGringotts.BridgeOutboundTransfer[]\",\"name\":\"outbounds\",\"type\":\"tuple[]\"}],\"internalType\":\"structGringotts.BridgeRequest\",\"name\":\"_params\",\"type\":\"tuple\"}],\"name\":\"bridge\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"messageIds\",\"type\":\"bytes32[]\"}],\"internalType\":\"structGringotts.BridgeResponse\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"endpoint\",\"outputs\":[{\"internalType\":\"contractILayerZeroEndpointV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amountUSDX\",\"type\":\"uint256\"}],\"internalType\":\"structGringotts.EstimateInboundTransfer\",\"name\":\"inbound\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"ChainId\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"executionGas\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"messageLength\",\"type\":\"uint16\"}],\"internalType\":\"structGringotts.EstimateOutboundTransfer[]\",\"name\":\"outbounds\",\"type\":\"tuple[]\"}],\"internalType\":\"structGringotts.EstimateRequest\",\"name\":\"_params\",\"type\":\"tuple\"}],\"name\":\"estimate\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"commissionUSDX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commissionDiscountUSDX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transferGasUSDX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transferGasDiscountUSDX\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"ChainId\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"executionGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transferGas\",\"type\":\"uint256\"}],\"internalType\":\"structGringotts.EstimateOutboundDetails[]\",\"name\":\"outboundDetails\",\"type\":\"tuple[]\"}],\"internalType\":\"structGringotts.EstimateResponse\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"srcEid\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structOrigin\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"}],\"name\":\"isComposeMsgSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"srcEid\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"internalType\":\"structOrigin\",\"name\":\"_origin\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"_guid\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"lzReceive\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"nextNonce\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oAppVersion\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"senderVersion\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"receiverVersion\",\"type\":\"uint64\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"eid\",\"type\":\"uint32\"}],\"name\":\"peers\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"peer\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_chainlinkPriceFeed\",\"type\":\"address\"}],\"name\":\"setChainlinkPriceFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"commissionMicroBPS\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"commissionDiscountBPS\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasDiscountBPS\",\"type\":\"uint32\"},{\"internalType\":\"GringottsAddress[]\",\"name\":\"stableCoins\",\"type\":\"bytes32[]\"}],\"internalType\":\"structConfig\",\"name\":\"_config\",\"type\":\"tuple\"}],\"name\":\"setConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegate\",\"type\":\"address\"}],\"name\":\"setDelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_eid\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"_peer\",\"type\":\"bytes32\"}],\"name\":\"setPeer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pythPriceFeed\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_pythPriceFeedId\",\"type\":\"bytes32\"}],\"name\":\"setPythPriceFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_winklinkPriceFeed\",\"type\":\"address\"}],\"name\":\"setWinklinkPriceFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"ChainId\",\"name\":\"_chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_header\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_header2\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"},{\"internalType\":\"uint128\",\"name\":\"_gas\",\"type\":\"uint128\"},{\"internalType\":\"bool\",\"name\":\"multi\",\"type\":\"bool\"}],\"name\":\"testSend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_dstEid\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"_m\",\"type\":\"string\"}],\"name\":\"testSend\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"unblockAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"ChainId\",\"name\":\"_chainID\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"ChainId\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"GringottsAddress\",\"name\":\"endpoint\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"lzEID\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"multiSend\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"baseGasEstimate\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"registerGasEstimate\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"completionGasEstimate\",\"type\":\"uint128\"}],\"internalType\":\"structPeer\",\"name\":\"_agent\",\"type\":\"tuple\"}],\"name\":\"updateAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// GringottsEVMABI is the input ABI used to generate the binding from.
// Deprecated: Use GringottsEVMMetaData.ABI instead.
var GringottsEVMABI = GringottsEVMMetaData.ABI

// GringottsEVM is an auto generated Go binding around an Ethereum contract.
type GringottsEVM struct {
	GringottsEVMCaller     // Read-only binding to the contract
	GringottsEVMTransactor // Write-only binding to the contract
	GringottsEVMFilterer   // Log filterer for contract events
}

// GringottsEVMCaller is an auto generated read-only Go binding around an Ethereum contract.
type GringottsEVMCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GringottsEVMTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GringottsEVMTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GringottsEVMFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GringottsEVMFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GringottsEVMSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GringottsEVMSession struct {
	Contract     *GringottsEVM     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Swap auth options to use throughout this session
}

// GringottsEVMCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GringottsEVMCallerSession struct {
	Contract *GringottsEVMCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GringottsEVMTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GringottsEVMTransactorSession struct {
	Contract     *GringottsEVMTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Swap auth options to use throughout this session
}

// GringottsEVMRaw is an auto generated low-level Go binding around an Ethereum contract.
type GringottsEVMRaw struct {
	Contract *GringottsEVM // Generic contract binding to access the raw methods on
}

// GringottsEVMCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GringottsEVMCallerRaw struct {
	Contract *GringottsEVMCaller // Generic read-only contract binding to access the raw methods on
}

// GringottsEVMTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GringottsEVMTransactorRaw struct {
	Contract *GringottsEVMTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGringottsEVM creates a new instance of GringottsEVM, bound to a specific deployed contract.
func NewGringottsEVM(address common.Address, backend bind.ContractBackend) (*GringottsEVM, error) {
	contract, err := bindGringottsEVM(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GringottsEVM{GringottsEVMCaller: GringottsEVMCaller{contract: contract}, GringottsEVMTransactor: GringottsEVMTransactor{contract: contract}, GringottsEVMFilterer: GringottsEVMFilterer{contract: contract}}, nil
}

// NewGringottsEVMCaller creates a new read-only instance of GringottsEVM, bound to a specific deployed contract.
func NewGringottsEVMCaller(address common.Address, caller bind.ContractCaller) (*GringottsEVMCaller, error) {
	contract, err := bindGringottsEVM(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMCaller{contract: contract}, nil
}

// NewGringottsEVMTransactor creates a new write-only instance of GringottsEVM, bound to a specific deployed contract.
func NewGringottsEVMTransactor(address common.Address, transactor bind.ContractTransactor) (*GringottsEVMTransactor, error) {
	contract, err := bindGringottsEVM(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMTransactor{contract: contract}, nil
}

// NewGringottsEVMFilterer creates a new log filterer instance of GringottsEVM, bound to a specific deployed contract.
func NewGringottsEVMFilterer(address common.Address, filterer bind.ContractFilterer) (*GringottsEVMFilterer, error) {
	contract, err := bindGringottsEVM(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMFilterer{contract: contract}, nil
}

// bindGringottsEVM binds a generic wrapper to an already deployed contract.
func bindGringottsEVM(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GringottsEVMABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GringottsEVM *GringottsEVMRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GringottsEVM.Contract.GringottsEVMCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GringottsEVM *GringottsEVMRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GringottsEVM.Contract.GringottsEVMTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GringottsEVM *GringottsEVMRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GringottsEVM.Contract.GringottsEVMTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GringottsEVM *GringottsEVMCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GringottsEVM.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GringottsEVM *GringottsEVMTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GringottsEVM.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GringottsEVM *GringottsEVMTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GringottsEVM.Contract.contract.Transact(opts, method, params...)
}

// AllowInitializePath is a free data retrieval call binding the contract method 0xff7bd03d.
//
// Solidity: function allowInitializePath((uint32,bytes32,uint64) origin) view returns(bool)
func (_GringottsEVM *GringottsEVMCaller) AllowInitializePath(opts *bind.CallOpts, origin Origin) (bool, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "allowInitializePath", origin)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowInitializePath is a free data retrieval call binding the contract method 0xff7bd03d.
//
// Solidity: function allowInitializePath((uint32,bytes32,uint64) origin) view returns(bool)
func (_GringottsEVM *GringottsEVMSession) AllowInitializePath(origin Origin) (bool, error) {
	return _GringottsEVM.Contract.AllowInitializePath(&_GringottsEVM.CallOpts, origin)
}

// AllowInitializePath is a free data retrieval call binding the contract method 0xff7bd03d.
//
// Solidity: function allowInitializePath((uint32,bytes32,uint64) origin) view returns(bool)
func (_GringottsEVM *GringottsEVMCallerSession) AllowInitializePath(origin Origin) (bool, error) {
	return _GringottsEVM.Contract.AllowInitializePath(&_GringottsEVM.CallOpts, origin)
}

// BalanceERC20 is a free data retrieval call binding the contract method 0x0f4ac71a.
//
// Solidity: function balanceERC20(address _token) view returns(uint256)
func (_GringottsEVM *GringottsEVMCaller) BalanceERC20(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "balanceERC20", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceERC20 is a free data retrieval call binding the contract method 0x0f4ac71a.
//
// Solidity: function balanceERC20(address _token) view returns(uint256)
func (_GringottsEVM *GringottsEVMSession) BalanceERC20(_token common.Address) (*big.Int, error) {
	return _GringottsEVM.Contract.BalanceERC20(&_GringottsEVM.CallOpts, _token)
}

// BalanceERC20 is a free data retrieval call binding the contract method 0x0f4ac71a.
//
// Solidity: function balanceERC20(address _token) view returns(uint256)
func (_GringottsEVM *GringottsEVMCallerSession) BalanceERC20(_token common.Address) (*big.Int, error) {
	return _GringottsEVM.Contract.BalanceERC20(&_GringottsEVM.CallOpts, _token)
}

// Endpoint is a free data retrieval call binding the contract method 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (_GringottsEVM *GringottsEVMCaller) Endpoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "endpoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Endpoint is a free data retrieval call binding the contract method 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (_GringottsEVM *GringottsEVMSession) Endpoint() (common.Address, error) {
	return _GringottsEVM.Contract.Endpoint(&_GringottsEVM.CallOpts)
}

// Endpoint is a free data retrieval call binding the contract method 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (_GringottsEVM *GringottsEVMCallerSession) Endpoint() (common.Address, error) {
	return _GringottsEVM.Contract.Endpoint(&_GringottsEVM.CallOpts)
}

// Estimate is a free data retrieval call binding the contract method 0x2a42ca36.
//
// Solidity: function estimate(((uint256),(uint8,uint256,uint16)[]) _params) view returns((uint256,uint256,uint256,uint256,(uint8,uint256,uint256)[]))
func (_GringottsEVM *GringottsEVMCaller) Estimate(opts *bind.CallOpts, _params GringottsEstimateRequest) (GringottsEstimateResponse, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "estimate", _params)

	if err != nil {
		return *new(GringottsEstimateResponse), err
	}

	out0 := *abi.ConvertType(out[0], new(GringottsEstimateResponse)).(*GringottsEstimateResponse)

	return out0, err

}

// Estimate is a free data retrieval call binding the contract method 0x2a42ca36.
//
// Solidity: function estimate(((uint256),(uint8,uint256,uint16)[]) _params) view returns((uint256,uint256,uint256,uint256,(uint8,uint256,uint256)[]))
func (_GringottsEVM *GringottsEVMSession) Estimate(_params GringottsEstimateRequest) (GringottsEstimateResponse, error) {
	return _GringottsEVM.Contract.Estimate(&_GringottsEVM.CallOpts, _params)
}

// Estimate is a free data retrieval call binding the contract method 0x2a42ca36.
//
// Solidity: function estimate(((uint256),(uint8,uint256,uint16)[]) _params) view returns((uint256,uint256,uint256,uint256,(uint8,uint256,uint256)[]))
func (_GringottsEVM *GringottsEVMCallerSession) Estimate(_params GringottsEstimateRequest) (GringottsEstimateResponse, error) {
	return _GringottsEVM.Contract.Estimate(&_GringottsEVM.CallOpts, _params)
}

// IsComposeMsgSender is a free data retrieval call binding the contract method 0x82413eac.
//
// Solidity: function isComposeMsgSender((uint32,bytes32,uint64) , bytes , address _sender) view returns(bool)
func (_GringottsEVM *GringottsEVMCaller) IsComposeMsgSender(opts *bind.CallOpts, arg0 Origin, arg1 []byte, _sender common.Address) (bool, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "isComposeMsgSender", arg0, arg1, _sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsComposeMsgSender is a free data retrieval call binding the contract method 0x82413eac.
//
// Solidity: function isComposeMsgSender((uint32,bytes32,uint64) , bytes , address _sender) view returns(bool)
func (_GringottsEVM *GringottsEVMSession) IsComposeMsgSender(arg0 Origin, arg1 []byte, _sender common.Address) (bool, error) {
	return _GringottsEVM.Contract.IsComposeMsgSender(&_GringottsEVM.CallOpts, arg0, arg1, _sender)
}

// IsComposeMsgSender is a free data retrieval call binding the contract method 0x82413eac.
//
// Solidity: function isComposeMsgSender((uint32,bytes32,uint64) , bytes , address _sender) view returns(bool)
func (_GringottsEVM *GringottsEVMCallerSession) IsComposeMsgSender(arg0 Origin, arg1 []byte, _sender common.Address) (bool, error) {
	return _GringottsEVM.Contract.IsComposeMsgSender(&_GringottsEVM.CallOpts, arg0, arg1, _sender)
}

// NextNonce is a free data retrieval call binding the contract method 0x7d25a05e.
//
// Solidity: function nextNonce(uint32 , bytes32 ) view returns(uint64 nonce)
func (_GringottsEVM *GringottsEVMCaller) NextNonce(opts *bind.CallOpts, arg0 uint32, arg1 [32]byte) (uint64, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "nextNonce", arg0, arg1)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// NextNonce is a free data retrieval call binding the contract method 0x7d25a05e.
//
// Solidity: function nextNonce(uint32 , bytes32 ) view returns(uint64 nonce)
func (_GringottsEVM *GringottsEVMSession) NextNonce(arg0 uint32, arg1 [32]byte) (uint64, error) {
	return _GringottsEVM.Contract.NextNonce(&_GringottsEVM.CallOpts, arg0, arg1)
}

// NextNonce is a free data retrieval call binding the contract method 0x7d25a05e.
//
// Solidity: function nextNonce(uint32 , bytes32 ) view returns(uint64 nonce)
func (_GringottsEVM *GringottsEVMCallerSession) NextNonce(arg0 uint32, arg1 [32]byte) (uint64, error) {
	return _GringottsEVM.Contract.NextNonce(&_GringottsEVM.CallOpts, arg0, arg1)
}

// OAppVersion is a free data retrieval call binding the contract method 0x17442b70.
//
// Solidity: function oAppVersion() pure returns(uint64 senderVersion, uint64 receiverVersion)
func (_GringottsEVM *GringottsEVMCaller) OAppVersion(opts *bind.CallOpts) (struct {
	SenderVersion   uint64
	ReceiverVersion uint64
}, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "oAppVersion")

	outstruct := new(struct {
		SenderVersion   uint64
		ReceiverVersion uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SenderVersion = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.ReceiverVersion = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// OAppVersion is a free data retrieval call binding the contract method 0x17442b70.
//
// Solidity: function oAppVersion() pure returns(uint64 senderVersion, uint64 receiverVersion)
func (_GringottsEVM *GringottsEVMSession) OAppVersion() (struct {
	SenderVersion   uint64
	ReceiverVersion uint64
}, error) {
	return _GringottsEVM.Contract.OAppVersion(&_GringottsEVM.CallOpts)
}

// OAppVersion is a free data retrieval call binding the contract method 0x17442b70.
//
// Solidity: function oAppVersion() pure returns(uint64 senderVersion, uint64 receiverVersion)
func (_GringottsEVM *GringottsEVMCallerSession) OAppVersion() (struct {
	SenderVersion   uint64
	ReceiverVersion uint64
}, error) {
	return _GringottsEVM.Contract.OAppVersion(&_GringottsEVM.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GringottsEVM *GringottsEVMCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GringottsEVM *GringottsEVMSession) Owner() (common.Address, error) {
	return _GringottsEVM.Contract.Owner(&_GringottsEVM.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GringottsEVM *GringottsEVMCallerSession) Owner() (common.Address, error) {
	return _GringottsEVM.Contract.Owner(&_GringottsEVM.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GringottsEVM *GringottsEVMCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GringottsEVM *GringottsEVMSession) Paused() (bool, error) {
	return _GringottsEVM.Contract.Paused(&_GringottsEVM.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GringottsEVM *GringottsEVMCallerSession) Paused() (bool, error) {
	return _GringottsEVM.Contract.Paused(&_GringottsEVM.CallOpts)
}

// Peers is a free data retrieval call binding the contract method 0xbb0b6a53.
//
// Solidity: function peers(uint32 eid) view returns(bytes32 peer)
func (_GringottsEVM *GringottsEVMCaller) Peers(opts *bind.CallOpts, eid uint32) ([32]byte, error) {
	var out []interface{}
	err := _GringottsEVM.contract.Call(opts, &out, "peers", eid)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Peers is a free data retrieval call binding the contract method 0xbb0b6a53.
//
// Solidity: function peers(uint32 eid) view returns(bytes32 peer)
func (_GringottsEVM *GringottsEVMSession) Peers(eid uint32) ([32]byte, error) {
	return _GringottsEVM.Contract.Peers(&_GringottsEVM.CallOpts, eid)
}

// Peers is a free data retrieval call binding the contract method 0xbb0b6a53.
//
// Solidity: function peers(uint32 eid) view returns(bytes32 peer)
func (_GringottsEVM *GringottsEVMCallerSession) Peers(eid uint32) ([32]byte, error) {
	return _GringottsEVM.Contract.Peers(&_GringottsEVM.CallOpts, eid)
}

// BlockAccount is a paid mutator transaction binding the contract method 0x7c0a893d.
//
// Solidity: function blockAccount(address account) returns()
func (_GringottsEVM *GringottsEVMTransactor) BlockAccount(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "blockAccount", account)
}

// BlockAccount is a paid mutator transaction binding the contract method 0x7c0a893d.
//
// Solidity: function blockAccount(address account) returns()
func (_GringottsEVM *GringottsEVMSession) BlockAccount(account common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.BlockAccount(&_GringottsEVM.TransactOpts, account)
}

// BlockAccount is a paid mutator transaction binding the contract method 0x7c0a893d.
//
// Solidity: function blockAccount(address account) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) BlockAccount(account common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.BlockAccount(&_GringottsEVM.TransactOpts, account)
}

// Bridge is a paid mutator transaction binding the contract method 0xfea18532.
//
// Solidity: function bridge(((uint256,(bytes32,uint256,bytes32,bytes,bytes32)[]),(uint8,uint256,bytes,(uint16)[])[]) _params) payable returns((bytes32[]))
func (_GringottsEVM *GringottsEVMTransactor) Bridge(opts *bind.TransactOpts, _params GringottsBridgeRequest) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "bridge", _params)
}

// Bridge is a paid mutator transaction binding the contract method 0xfea18532.
//
// Solidity: function bridge(((uint256,(bytes32,uint256,bytes32,bytes,bytes32)[]),(uint8,uint256,bytes,(uint16)[])[]) _params) payable returns((bytes32[]))
func (_GringottsEVM *GringottsEVMSession) Bridge(_params GringottsBridgeRequest) (*types.Transaction, error) {
	return _GringottsEVM.Contract.Bridge(&_GringottsEVM.TransactOpts, _params)
}

// Bridge is a paid mutator transaction binding the contract method 0xfea18532.
//
// Solidity: function bridge(((uint256,(bytes32,uint256,bytes32,bytes,bytes32)[]),(uint8,uint256,bytes,(uint16)[])[]) _params) payable returns((bytes32[]))
func (_GringottsEVM *GringottsEVMTransactorSession) Bridge(_params GringottsBridgeRequest) (*types.Transaction, error) {
	return _GringottsEVM.Contract.Bridge(&_GringottsEVM.TransactOpts, _params)
}

// LzReceive is a paid mutator transaction binding the contract method 0x13137d65.
//
// Solidity: function lzReceive((uint32,bytes32,uint64) _origin, bytes32 _guid, bytes _message, address _executor, bytes _extraData) payable returns()
func (_GringottsEVM *GringottsEVMTransactor) LzReceive(opts *bind.TransactOpts, _origin Origin, _guid [32]byte, _message []byte, _executor common.Address, _extraData []byte) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "lzReceive", _origin, _guid, _message, _executor, _extraData)
}

// LzReceive is a paid mutator transaction binding the contract method 0x13137d65.
//
// Solidity: function lzReceive((uint32,bytes32,uint64) _origin, bytes32 _guid, bytes _message, address _executor, bytes _extraData) payable returns()
func (_GringottsEVM *GringottsEVMSession) LzReceive(_origin Origin, _guid [32]byte, _message []byte, _executor common.Address, _extraData []byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.LzReceive(&_GringottsEVM.TransactOpts, _origin, _guid, _message, _executor, _extraData)
}

// LzReceive is a paid mutator transaction binding the contract method 0x13137d65.
//
// Solidity: function lzReceive((uint32,bytes32,uint64) _origin, bytes32 _guid, bytes _message, address _executor, bytes _extraData) payable returns()
func (_GringottsEVM *GringottsEVMTransactorSession) LzReceive(_origin Origin, _guid [32]byte, _message []byte, _executor common.Address, _extraData []byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.LzReceive(&_GringottsEVM.TransactOpts, _origin, _guid, _message, _executor, _extraData)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GringottsEVM *GringottsEVMTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GringottsEVM *GringottsEVMSession) Pause() (*types.Transaction, error) {
	return _GringottsEVM.Contract.Pause(&_GringottsEVM.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GringottsEVM *GringottsEVMTransactorSession) Pause() (*types.Transaction, error) {
	return _GringottsEVM.Contract.Pause(&_GringottsEVM.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GringottsEVM *GringottsEVMTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GringottsEVM *GringottsEVMSession) RenounceOwnership() (*types.Transaction, error) {
	return _GringottsEVM.Contract.RenounceOwnership(&_GringottsEVM.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GringottsEVM *GringottsEVMTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GringottsEVM.Contract.RenounceOwnership(&_GringottsEVM.TransactOpts)
}

// SetChainlinkPriceFeed is a paid mutator transaction binding the contract method 0xf88f964a.
//
// Solidity: function setChainlinkPriceFeed(address _chainlinkPriceFeed) returns()
func (_GringottsEVM *GringottsEVMTransactor) SetChainlinkPriceFeed(opts *bind.TransactOpts, _chainlinkPriceFeed common.Address) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "setChainlinkPriceFeed", _chainlinkPriceFeed)
}

// SetChainlinkPriceFeed is a paid mutator transaction binding the contract method 0xf88f964a.
//
// Solidity: function setChainlinkPriceFeed(address _chainlinkPriceFeed) returns()
func (_GringottsEVM *GringottsEVMSession) SetChainlinkPriceFeed(_chainlinkPriceFeed common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetChainlinkPriceFeed(&_GringottsEVM.TransactOpts, _chainlinkPriceFeed)
}

// SetChainlinkPriceFeed is a paid mutator transaction binding the contract method 0xf88f964a.
//
// Solidity: function setChainlinkPriceFeed(address _chainlinkPriceFeed) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) SetChainlinkPriceFeed(_chainlinkPriceFeed common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetChainlinkPriceFeed(&_GringottsEVM.TransactOpts, _chainlinkPriceFeed)
}

// SetConfig is a paid mutator transaction binding the contract method 0x77bff184.
//
// Solidity: function setConfig((uint32,uint32,uint32,bytes32[]) _config) returns()
func (_GringottsEVM *GringottsEVMTransactor) SetConfig(opts *bind.TransactOpts, _config Config) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "setConfig", _config)
}

// SetConfig is a paid mutator transaction binding the contract method 0x77bff184.
//
// Solidity: function setConfig((uint32,uint32,uint32,bytes32[]) _config) returns()
func (_GringottsEVM *GringottsEVMSession) SetConfig(_config Config) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetConfig(&_GringottsEVM.TransactOpts, _config)
}

// SetConfig is a paid mutator transaction binding the contract method 0x77bff184.
//
// Solidity: function setConfig((uint32,uint32,uint32,bytes32[]) _config) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) SetConfig(_config Config) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetConfig(&_GringottsEVM.TransactOpts, _config)
}

// SetDelegate is a paid mutator transaction binding the contract method 0xca5eb5e1.
//
// Solidity: function setDelegate(address _delegate) returns()
func (_GringottsEVM *GringottsEVMTransactor) SetDelegate(opts *bind.TransactOpts, _delegate common.Address) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "setDelegate", _delegate)
}

// SetDelegate is a paid mutator transaction binding the contract method 0xca5eb5e1.
//
// Solidity: function setDelegate(address _delegate) returns()
func (_GringottsEVM *GringottsEVMSession) SetDelegate(_delegate common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetDelegate(&_GringottsEVM.TransactOpts, _delegate)
}

// SetDelegate is a paid mutator transaction binding the contract method 0xca5eb5e1.
//
// Solidity: function setDelegate(address _delegate) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) SetDelegate(_delegate common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetDelegate(&_GringottsEVM.TransactOpts, _delegate)
}

// SetPeer is a paid mutator transaction binding the contract method 0x3400288b.
//
// Solidity: function setPeer(uint32 _eid, bytes32 _peer) returns()
func (_GringottsEVM *GringottsEVMTransactor) SetPeer(opts *bind.TransactOpts, _eid uint32, _peer [32]byte) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "setPeer", _eid, _peer)
}

// SetPeer is a paid mutator transaction binding the contract method 0x3400288b.
//
// Solidity: function setPeer(uint32 _eid, bytes32 _peer) returns()
func (_GringottsEVM *GringottsEVMSession) SetPeer(_eid uint32, _peer [32]byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetPeer(&_GringottsEVM.TransactOpts, _eid, _peer)
}

// SetPeer is a paid mutator transaction binding the contract method 0x3400288b.
//
// Solidity: function setPeer(uint32 _eid, bytes32 _peer) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) SetPeer(_eid uint32, _peer [32]byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetPeer(&_GringottsEVM.TransactOpts, _eid, _peer)
}

// SetPythPriceFeed is a paid mutator transaction binding the contract method 0x28615ad2.
//
// Solidity: function setPythPriceFeed(address _pythPriceFeed, bytes32 _pythPriceFeedId) returns()
func (_GringottsEVM *GringottsEVMTransactor) SetPythPriceFeed(opts *bind.TransactOpts, _pythPriceFeed common.Address, _pythPriceFeedId [32]byte) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "setPythPriceFeed", _pythPriceFeed, _pythPriceFeedId)
}

// SetPythPriceFeed is a paid mutator transaction binding the contract method 0x28615ad2.
//
// Solidity: function setPythPriceFeed(address _pythPriceFeed, bytes32 _pythPriceFeedId) returns()
func (_GringottsEVM *GringottsEVMSession) SetPythPriceFeed(_pythPriceFeed common.Address, _pythPriceFeedId [32]byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetPythPriceFeed(&_GringottsEVM.TransactOpts, _pythPriceFeed, _pythPriceFeedId)
}

// SetPythPriceFeed is a paid mutator transaction binding the contract method 0x28615ad2.
//
// Solidity: function setPythPriceFeed(address _pythPriceFeed, bytes32 _pythPriceFeedId) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) SetPythPriceFeed(_pythPriceFeed common.Address, _pythPriceFeedId [32]byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetPythPriceFeed(&_GringottsEVM.TransactOpts, _pythPriceFeed, _pythPriceFeedId)
}

// SetWinklinkPriceFeed is a paid mutator transaction binding the contract method 0x91a2827d.
//
// Solidity: function setWinklinkPriceFeed(address _winklinkPriceFeed) returns()
func (_GringottsEVM *GringottsEVMTransactor) SetWinklinkPriceFeed(opts *bind.TransactOpts, _winklinkPriceFeed common.Address) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "setWinklinkPriceFeed", _winklinkPriceFeed)
}

// SetWinklinkPriceFeed is a paid mutator transaction binding the contract method 0x91a2827d.
//
// Solidity: function setWinklinkPriceFeed(address _winklinkPriceFeed) returns()
func (_GringottsEVM *GringottsEVMSession) SetWinklinkPriceFeed(_winklinkPriceFeed common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetWinklinkPriceFeed(&_GringottsEVM.TransactOpts, _winklinkPriceFeed)
}

// SetWinklinkPriceFeed is a paid mutator transaction binding the contract method 0x91a2827d.
//
// Solidity: function setWinklinkPriceFeed(address _winklinkPriceFeed) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) SetWinklinkPriceFeed(_winklinkPriceFeed common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.SetWinklinkPriceFeed(&_GringottsEVM.TransactOpts, _winklinkPriceFeed)
}

// TestSend is a paid mutator transaction binding the contract method 0x055e877b.
//
// Solidity: function testSend(uint8 _chainID, uint8 _header, uint8 _header2, bytes _message, uint128 _gas, bool multi) returns()
func (_GringottsEVM *GringottsEVMTransactor) TestSend(opts *bind.TransactOpts, _chainID uint8, _header uint8, _header2 uint8, _message []byte, _gas *big.Int, multi bool) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "testSend", _chainID, _header, _header2, _message, _gas, multi)
}

// TestSend is a paid mutator transaction binding the contract method 0x055e877b.
//
// Solidity: function testSend(uint8 _chainID, uint8 _header, uint8 _header2, bytes _message, uint128 _gas, bool multi) returns()
func (_GringottsEVM *GringottsEVMSession) TestSend(_chainID uint8, _header uint8, _header2 uint8, _message []byte, _gas *big.Int, multi bool) (*types.Transaction, error) {
	return _GringottsEVM.Contract.TestSend(&_GringottsEVM.TransactOpts, _chainID, _header, _header2, _message, _gas, multi)
}

// TestSend is a paid mutator transaction binding the contract method 0x055e877b.
//
// Solidity: function testSend(uint8 _chainID, uint8 _header, uint8 _header2, bytes _message, uint128 _gas, bool multi) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) TestSend(_chainID uint8, _header uint8, _header2 uint8, _message []byte, _gas *big.Int, multi bool) (*types.Transaction, error) {
	return _GringottsEVM.Contract.TestSend(&_GringottsEVM.TransactOpts, _chainID, _header, _header2, _message, _gas, multi)
}

// TestSend0 is a paid mutator transaction binding the contract method 0x0bcd6333.
//
// Solidity: function testSend(uint32 _dstEid, string _m) returns(string)
func (_GringottsEVM *GringottsEVMTransactor) TestSend0(opts *bind.TransactOpts, _dstEid uint32, _m string) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "testSend0", _dstEid, _m)
}

// TestSend0 is a paid mutator transaction binding the contract method 0x0bcd6333.
//
// Solidity: function testSend(uint32 _dstEid, string _m) returns(string)
func (_GringottsEVM *GringottsEVMSession) TestSend0(_dstEid uint32, _m string) (*types.Transaction, error) {
	return _GringottsEVM.Contract.TestSend0(&_GringottsEVM.TransactOpts, _dstEid, _m)
}

// TestSend0 is a paid mutator transaction binding the contract method 0x0bcd6333.
//
// Solidity: function testSend(uint32 _dstEid, string _m) returns(string)
func (_GringottsEVM *GringottsEVMTransactorSession) TestSend0(_dstEid uint32, _m string) (*types.Transaction, error) {
	return _GringottsEVM.Contract.TestSend0(&_GringottsEVM.TransactOpts, _dstEid, _m)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GringottsEVM *GringottsEVMTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GringottsEVM *GringottsEVMSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.TransferOwnership(&_GringottsEVM.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.TransferOwnership(&_GringottsEVM.TransactOpts, newOwner)
}

// UnblockAccount is a paid mutator transaction binding the contract method 0x4d78fdc6.
//
// Solidity: function unblockAccount(address account) returns()
func (_GringottsEVM *GringottsEVMTransactor) UnblockAccount(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "unblockAccount", account)
}

// UnblockAccount is a paid mutator transaction binding the contract method 0x4d78fdc6.
//
// Solidity: function unblockAccount(address account) returns()
func (_GringottsEVM *GringottsEVMSession) UnblockAccount(account common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.UnblockAccount(&_GringottsEVM.TransactOpts, account)
}

// UnblockAccount is a paid mutator transaction binding the contract method 0x4d78fdc6.
//
// Solidity: function unblockAccount(address account) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) UnblockAccount(account common.Address) (*types.Transaction, error) {
	return _GringottsEVM.Contract.UnblockAccount(&_GringottsEVM.TransactOpts, account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GringottsEVM *GringottsEVMTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GringottsEVM *GringottsEVMSession) Unpause() (*types.Transaction, error) {
	return _GringottsEVM.Contract.Unpause(&_GringottsEVM.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GringottsEVM *GringottsEVMTransactorSession) Unpause() (*types.Transaction, error) {
	return _GringottsEVM.Contract.Unpause(&_GringottsEVM.TransactOpts)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xd3dcbdab.
//
// Solidity: function updateAgent(uint8 _chainID, (uint8,bytes32,uint32,bool,uint128,uint128,uint128) _agent) returns()
func (_GringottsEVM *GringottsEVMTransactor) UpdateAgent(opts *bind.TransactOpts, _chainID uint8, _agent Peer) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "updateAgent", _chainID, _agent)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xd3dcbdab.
//
// Solidity: function updateAgent(uint8 _chainID, (uint8,bytes32,uint32,bool,uint128,uint128,uint128) _agent) returns()
func (_GringottsEVM *GringottsEVMSession) UpdateAgent(_chainID uint8, _agent Peer) (*types.Transaction, error) {
	return _GringottsEVM.Contract.UpdateAgent(&_GringottsEVM.TransactOpts, _chainID, _agent)
}

// UpdateAgent is a paid mutator transaction binding the contract method 0xd3dcbdab.
//
// Solidity: function updateAgent(uint8 _chainID, (uint8,bytes32,uint32,bool,uint128,uint128,uint128) _agent) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) UpdateAgent(_chainID uint8, _agent Peer) (*types.Transaction, error) {
	return _GringottsEVM.Contract.UpdateAgent(&_GringottsEVM.TransactOpts, _chainID, _agent)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_GringottsEVM *GringottsEVMTransactor) WithdrawERC20(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "withdrawERC20", _token, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_GringottsEVM *GringottsEVMSession) WithdrawERC20(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _GringottsEVM.Contract.WithdrawERC20(&_GringottsEVM.TransactOpts, _token, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) WithdrawERC20(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _GringottsEVM.Contract.WithdrawERC20(&_GringottsEVM.TransactOpts, _token, _amount)
}

// WithdrawNative is a paid mutator transaction binding the contract method 0x84276d81.
//
// Solidity: function withdrawNative(uint256 _amount) returns()
func (_GringottsEVM *GringottsEVMTransactor) WithdrawNative(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _GringottsEVM.contract.Transact(opts, "withdrawNative", _amount)
}

// WithdrawNative is a paid mutator transaction binding the contract method 0x84276d81.
//
// Solidity: function withdrawNative(uint256 _amount) returns()
func (_GringottsEVM *GringottsEVMSession) WithdrawNative(_amount *big.Int) (*types.Transaction, error) {
	return _GringottsEVM.Contract.WithdrawNative(&_GringottsEVM.TransactOpts, _amount)
}

// WithdrawNative is a paid mutator transaction binding the contract method 0x84276d81.
//
// Solidity: function withdrawNative(uint256 _amount) returns()
func (_GringottsEVM *GringottsEVMTransactorSession) WithdrawNative(_amount *big.Int) (*types.Transaction, error) {
	return _GringottsEVM.Contract.WithdrawNative(&_GringottsEVM.TransactOpts, _amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_GringottsEVM *GringottsEVMTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _GringottsEVM.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_GringottsEVM *GringottsEVMSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.Fallback(&_GringottsEVM.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_GringottsEVM *GringottsEVMTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _GringottsEVM.Contract.Fallback(&_GringottsEVM.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GringottsEVM *GringottsEVMTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GringottsEVM.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GringottsEVM *GringottsEVMSession) Receive() (*types.Transaction, error) {
	return _GringottsEVM.Contract.Receive(&_GringottsEVM.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GringottsEVM *GringottsEVMTransactorSession) Receive() (*types.Transaction, error) {
	return _GringottsEVM.Contract.Receive(&_GringottsEVM.TransactOpts)
}

// GringottsEVMBlockEventIterator is returned from FilterBlockEvent and is used to iterate over the raw logs and unpacked data for BlockEvent events raised by the GringottsEVM contract.
type GringottsEVMBlockEventIterator struct {
	Event *GringottsEVMBlockEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMBlockEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMBlockEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMBlockEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMBlockEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMBlockEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMBlockEvent represents a BlockEvent event raised by the GringottsEVM contract.
type GringottsEVMBlockEvent struct {
	Account common.Address
	Blocked bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBlockEvent is a free log retrieval operation binding the contract event 0x426d2cfff938548ecd1197f0a17954b1d866f471f9609a0b32b803f1021f4a50.
//
// Solidity: event BlockEvent(address indexed account, bool blocked)
func (_GringottsEVM *GringottsEVMFilterer) FilterBlockEvent(opts *bind.FilterOpts, account []common.Address) (*GringottsEVMBlockEventIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "BlockEvent", accountRule)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMBlockEventIterator{contract: _GringottsEVM.contract, event: "BlockEvent", logs: logs, sub: sub}, nil
}

// WatchBlockEvent is a free log subscription operation binding the contract event 0x426d2cfff938548ecd1197f0a17954b1d866f471f9609a0b32b803f1021f4a50.
//
// Solidity: event BlockEvent(address indexed account, bool blocked)
func (_GringottsEVM *GringottsEVMFilterer) WatchBlockEvent(opts *bind.WatchOpts, sink chan<- *GringottsEVMBlockEvent, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "BlockEvent", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMBlockEvent)
				if err := _GringottsEVM.contract.UnpackLog(event, "BlockEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBlockEvent is a log parse operation binding the contract event 0x426d2cfff938548ecd1197f0a17954b1d866f471f9609a0b32b803f1021f4a50.
//
// Solidity: event BlockEvent(address indexed account, bool blocked)
func (_GringottsEVM *GringottsEVMFilterer) ParseBlockEvent(log types.Log) (*GringottsEVMBlockEvent, error) {
	event := new(GringottsEVMBlockEvent)
	if err := _GringottsEVM.contract.UnpackLog(event, "BlockEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GringottsEVM contract.
type GringottsEVMOwnershipTransferredIterator struct {
	Event *GringottsEVMOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMOwnershipTransferred represents a OwnershipTransferred event raised by the GringottsEVM contract.
type GringottsEVMOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GringottsEVM *GringottsEVMFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GringottsEVMOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMOwnershipTransferredIterator{contract: _GringottsEVM.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GringottsEVM *GringottsEVMFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GringottsEVMOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMOwnershipTransferred)
				if err := _GringottsEVM.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GringottsEVM *GringottsEVMFilterer) ParseOwnershipTransferred(log types.Log) (*GringottsEVMOwnershipTransferred, error) {
	event := new(GringottsEVMOwnershipTransferred)
	if err := _GringottsEVM.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GringottsEVM contract.
type GringottsEVMPausedIterator struct {
	Event *GringottsEVMPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMPaused represents a Paused event raised by the GringottsEVM contract.
type GringottsEVMPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GringottsEVM *GringottsEVMFilterer) FilterPaused(opts *bind.FilterOpts) (*GringottsEVMPausedIterator, error) {

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GringottsEVMPausedIterator{contract: _GringottsEVM.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GringottsEVM *GringottsEVMFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GringottsEVMPaused) (event.Subscription, error) {

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMPaused)
				if err := _GringottsEVM.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GringottsEVM *GringottsEVMFilterer) ParsePaused(log types.Log) (*GringottsEVMPaused, error) {
	event := new(GringottsEVMPaused)
	if err := _GringottsEVM.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMPeerSetIterator is returned from FilterPeerSet and is used to iterate over the raw logs and unpacked data for PeerSet events raised by the GringottsEVM contract.
type GringottsEVMPeerSetIterator struct {
	Event *GringottsEVMPeerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMPeerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMPeerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMPeerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMPeerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMPeerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMPeerSet represents a PeerSet event raised by the GringottsEVM contract.
type GringottsEVMPeerSet struct {
	Eid  uint32
	Peer [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPeerSet is a free log retrieval operation binding the contract event 0x238399d427b947898edb290f5ff0f9109849b1c3ba196a42e35f00c50a54b98b.
//
// Solidity: event PeerSet(uint32 eid, bytes32 peer)
func (_GringottsEVM *GringottsEVMFilterer) FilterPeerSet(opts *bind.FilterOpts) (*GringottsEVMPeerSetIterator, error) {

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "PeerSet")
	if err != nil {
		return nil, err
	}
	return &GringottsEVMPeerSetIterator{contract: _GringottsEVM.contract, event: "PeerSet", logs: logs, sub: sub}, nil
}

// WatchPeerSet is a free log subscription operation binding the contract event 0x238399d427b947898edb290f5ff0f9109849b1c3ba196a42e35f00c50a54b98b.
//
// Solidity: event PeerSet(uint32 eid, bytes32 peer)
func (_GringottsEVM *GringottsEVMFilterer) WatchPeerSet(opts *bind.WatchOpts, sink chan<- *GringottsEVMPeerSet) (event.Subscription, error) {

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "PeerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMPeerSet)
				if err := _GringottsEVM.contract.UnpackLog(event, "PeerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePeerSet is a log parse operation binding the contract event 0x238399d427b947898edb290f5ff0f9109849b1c3ba196a42e35f00c50a54b98b.
//
// Solidity: event PeerSet(uint32 eid, bytes32 peer)
func (_GringottsEVM *GringottsEVMFilterer) ParsePeerSet(log types.Log) (*GringottsEVMPeerSet, error) {
	event := new(GringottsEVMPeerSet)
	if err := _GringottsEVM.contract.UnpackLog(event, "PeerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMReceiveChainTransferEventIterator is returned from FilterReceiveChainTransferEvent and is used to iterate over the raw logs and unpacked data for ReceiveChainTransferEvent events raised by the GringottsEVM contract.
type GringottsEVMReceiveChainTransferEventIterator struct {
	Event *GringottsEVMReceiveChainTransferEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMReceiveChainTransferEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMReceiveChainTransferEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMReceiveChainTransferEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMReceiveChainTransferEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMReceiveChainTransferEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMReceiveChainTransferEvent represents a ReceiveChainTransferEvent event raised by the GringottsEVM contract.
type GringottsEVMReceiveChainTransferEvent struct {
	ChainId    uint8
	MessageId  common.Hash
	Asset      common.Address
	Recipient  common.Address
	AmountUSDX *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReceiveChainTransferEvent is a free log retrieval operation binding the contract event 0xd57fd8c519641fd0d48d398085ac7e409f2e00827088cfe317a3802bde2b674e.
//
// Solidity: event ReceiveChainTransferEvent(uint8 chainId, string indexed messageId, address asset, address recipient, uint256 amountUSDX)
func (_GringottsEVM *GringottsEVMFilterer) FilterReceiveChainTransferEvent(opts *bind.FilterOpts, messageId []string) (*GringottsEVMReceiveChainTransferEventIterator, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "ReceiveChainTransferEvent", messageIdRule)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMReceiveChainTransferEventIterator{contract: _GringottsEVM.contract, event: "ReceiveChainTransferEvent", logs: logs, sub: sub}, nil
}

// WatchReceiveChainTransferEvent is a free log subscription operation binding the contract event 0xd57fd8c519641fd0d48d398085ac7e409f2e00827088cfe317a3802bde2b674e.
//
// Solidity: event ReceiveChainTransferEvent(uint8 chainId, string indexed messageId, address asset, address recipient, uint256 amountUSDX)
func (_GringottsEVM *GringottsEVMFilterer) WatchReceiveChainTransferEvent(opts *bind.WatchOpts, sink chan<- *GringottsEVMReceiveChainTransferEvent, messageId []string) (event.Subscription, error) {

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "ReceiveChainTransferEvent", messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMReceiveChainTransferEvent)
				if err := _GringottsEVM.contract.UnpackLog(event, "ReceiveChainTransferEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReceiveChainTransferEvent is a log parse operation binding the contract event 0xd57fd8c519641fd0d48d398085ac7e409f2e00827088cfe317a3802bde2b674e.
//
// Solidity: event ReceiveChainTransferEvent(uint8 chainId, string indexed messageId, address asset, address recipient, uint256 amountUSDX)
func (_GringottsEVM *GringottsEVMFilterer) ParseReceiveChainTransferEvent(log types.Log) (*GringottsEVMReceiveChainTransferEvent, error) {
	event := new(GringottsEVMReceiveChainTransferEvent)
	if err := _GringottsEVM.contract.UnpackLog(event, "ReceiveChainTransferEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMSendChainTransferEventIterator is returned from FilterSendChainTransferEvent and is used to iterate over the raw logs and unpacked data for SendChainTransferEvent events raised by the GringottsEVM contract.
type GringottsEVMSendChainTransferEventIterator struct {
	Event *GringottsEVMSendChainTransferEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMSendChainTransferEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMSendChainTransferEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMSendChainTransferEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMSendChainTransferEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMSendChainTransferEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMSendChainTransferEvent represents a SendChainTransferEvent event raised by the GringottsEVM contract.
type GringottsEVMSendChainTransferEvent struct {
	User       common.Address
	ChainId    uint8
	MessageId  common.Hash
	AmountUSDX *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSendChainTransferEvent is a free log retrieval operation binding the contract event 0xe9661c3d243e1a480444bebfeeacf4307666da73067d3f37280fba3bfc299cb8.
//
// Solidity: event SendChainTransferEvent(address indexed user, uint8 chainId, string indexed messageId, uint256 amountUSDX)
func (_GringottsEVM *GringottsEVMFilterer) FilterSendChainTransferEvent(opts *bind.FilterOpts, user []common.Address, messageId []string) (*GringottsEVMSendChainTransferEventIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "SendChainTransferEvent", userRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return &GringottsEVMSendChainTransferEventIterator{contract: _GringottsEVM.contract, event: "SendChainTransferEvent", logs: logs, sub: sub}, nil
}

// WatchSendChainTransferEvent is a free log subscription operation binding the contract event 0xe9661c3d243e1a480444bebfeeacf4307666da73067d3f37280fba3bfc299cb8.
//
// Solidity: event SendChainTransferEvent(address indexed user, uint8 chainId, string indexed messageId, uint256 amountUSDX)
func (_GringottsEVM *GringottsEVMFilterer) WatchSendChainTransferEvent(opts *bind.WatchOpts, sink chan<- *GringottsEVMSendChainTransferEvent, user []common.Address, messageId []string) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var messageIdRule []interface{}
	for _, messageIdItem := range messageId {
		messageIdRule = append(messageIdRule, messageIdItem)
	}

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "SendChainTransferEvent", userRule, messageIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMSendChainTransferEvent)
				if err := _GringottsEVM.contract.UnpackLog(event, "SendChainTransferEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSendChainTransferEvent is a log parse operation binding the contract event 0xe9661c3d243e1a480444bebfeeacf4307666da73067d3f37280fba3bfc299cb8.
//
// Solidity: event SendChainTransferEvent(address indexed user, uint8 chainId, string indexed messageId, uint256 amountUSDX)
func (_GringottsEVM *GringottsEVMFilterer) ParseSendChainTransferEvent(log types.Log) (*GringottsEVMSendChainTransferEvent, error) {
	event := new(GringottsEVMSendChainTransferEvent)
	if err := _GringottsEVM.contract.UnpackLog(event, "SendChainTransferEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMTestMessageIterator is returned from FilterTestMessage and is used to iterate over the raw logs and unpacked data for TestMessage events raised by the GringottsEVM contract.
type GringottsEVMTestMessageIterator struct {
	Event *GringottsEVMTestMessage // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMTestMessageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMTestMessage)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMTestMessage)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMTestMessageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMTestMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMTestMessage represents a TestMessage event raised by the GringottsEVM contract.
type GringottsEVMTestMessage struct {
	MessageId [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTestMessage is a free log retrieval operation binding the contract event 0x9ce6f8fd843514c8d475fb84af69f46569a647c2ba4c922d26ab28f04ad1ec58.
//
// Solidity: event TestMessage(bytes32 messageId)
func (_GringottsEVM *GringottsEVMFilterer) FilterTestMessage(opts *bind.FilterOpts) (*GringottsEVMTestMessageIterator, error) {

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "TestMessage")
	if err != nil {
		return nil, err
	}
	return &GringottsEVMTestMessageIterator{contract: _GringottsEVM.contract, event: "TestMessage", logs: logs, sub: sub}, nil
}

// WatchTestMessage is a free log subscription operation binding the contract event 0x9ce6f8fd843514c8d475fb84af69f46569a647c2ba4c922d26ab28f04ad1ec58.
//
// Solidity: event TestMessage(bytes32 messageId)
func (_GringottsEVM *GringottsEVMFilterer) WatchTestMessage(opts *bind.WatchOpts, sink chan<- *GringottsEVMTestMessage) (event.Subscription, error) {

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "TestMessage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMTestMessage)
				if err := _GringottsEVM.contract.UnpackLog(event, "TestMessage", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTestMessage is a log parse operation binding the contract event 0x9ce6f8fd843514c8d475fb84af69f46569a647c2ba4c922d26ab28f04ad1ec58.
//
// Solidity: event TestMessage(bytes32 messageId)
func (_GringottsEVM *GringottsEVMFilterer) ParseTestMessage(log types.Log) (*GringottsEVMTestMessage, error) {
	event := new(GringottsEVMTestMessage)
	if err := _GringottsEVM.contract.UnpackLog(event, "TestMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GringottsEVMUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GringottsEVM contract.
type GringottsEVMUnpausedIterator struct {
	Event *GringottsEVMUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GringottsEVMUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GringottsEVMUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GringottsEVMUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GringottsEVMUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GringottsEVMUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GringottsEVMUnpaused represents a Unpaused event raised by the GringottsEVM contract.
type GringottsEVMUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GringottsEVM *GringottsEVMFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GringottsEVMUnpausedIterator, error) {

	logs, sub, err := _GringottsEVM.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GringottsEVMUnpausedIterator{contract: _GringottsEVM.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GringottsEVM *GringottsEVMFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GringottsEVMUnpaused) (event.Subscription, error) {

	logs, sub, err := _GringottsEVM.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GringottsEVMUnpaused)
				if err := _GringottsEVM.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GringottsEVM *GringottsEVMFilterer) ParseUnpaused(log types.Log) (*GringottsEVMUnpaused, error) {
	event := new(GringottsEVMUnpaused)
	if err := _GringottsEVM.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
