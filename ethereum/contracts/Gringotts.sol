// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Blockable} from "./utils/Blockable.sol";
import {MessageTypeChainTransfer} from "./Message.sol";
import {ChainTransfer, ChainTransferItem, ChainTransferLibrary} from "./ChainTransfer.sol";
import {Config, GringottsAddress, ChainId} from "./Types.sol";
import {Peer, PeerLibrary} from "./Peer.sol";
import {Market} from "./Market.sol";
import {Vault} from "./Vault.sol";
import {AddressUtils, MathUtils} from "./utils/Utils.sol";
import {LayerZeroBridge} from "./LayerZeroMessenger.sol";
import {SendChainTransferEvent, ReceiveChainTransferEvent} from "./Events.sol";

import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {Address} from "@openzeppelin/contracts/utils/Address.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract Gringotts is Ownable, Blockable, LayerZeroBridge, Vault, Market, Pausable {
    uint8 private constant CHAIN_TRANSFER_DECIMALS = 6;
    uint8 private constant MAX_TRANSFER = 8;
    uint16 private constant MAX_DEX_COMMAND = 2 * 1024;
    uint8 private constant SWAP_GAS_USAGE_PERCENTAGE = 80;

    ChainId private chainId;
    Config private config;
    mapping(ChainId => Peer) private gringottsAgents;

    constructor(
        ChainId _chainId,
        uint8 _networkDecimals,
        address _lzEndpoint
    )
    Ownable(msg.sender)
    LayerZeroBridge(_lzEndpoint, msg.sender)
    Market(_networkDecimals)
    {
        chainId = _chainId;
    }

    struct Swap {
        GringottsAddress executor;
        bytes command;
        bytes metadata;
        GringottsAddress stableToken;
    }

    struct BridgeInboundTransferItem {
        GringottsAddress asset;
        uint256 amount;
        Swap swap;
    }

    struct BridgeInboundTransfer {
        uint256 amountUSDX;
        BridgeInboundTransferItem[] items;
    }

    struct BridgeOutboundTransferItem {
        GringottsAddress asset;
        GringottsAddress recipient;
        uint256 executionGas;
        uint16 distributionBP;
        Swap swap;
    }

    struct BridgeOutboundTransfer {
        ChainId chainId;
        bytes metadata;
        BridgeOutboundTransferItem[] items;
    }

    struct BridgeRequest {
        BridgeInboundTransfer inbound;
        BridgeOutboundTransfer[] outbounds;
    }

    struct BridgeResponse {
        string[] messageIds;
    }

    function bridge(
        BridgeRequest memory _params
    ) external onlyNotBlocked whenNotPaused payable returns (BridgeResponse memory) {
        validateBridge(_params);

        Peer memory self = gringottsAgents[chainId];

        uint256 netUSDX = 0;

        /*********** [Inbound transaction] ***********/
        for (uint256 i = 0; i < _params.inbound.items.length; i++) {
            BridgeInboundTransferItem memory item = _params.inbound.items[i];
            address inTransferAddress = AddressUtils.bytes32ToAddress(item.asset);

            if (PeerLibrary.hasStableCoin(self, item.asset)) {
                IERC20Metadata stableToken = IERC20Metadata(inTransferAddress);
                SafeERC20.safeTransferFrom(stableToken, msg.sender, address(this), item.amount);

                netUSDX = netUSDX + MathUtils.changeDecimals(
                    item.amount,
                    stableToken.decimals(),
                    CHAIN_TRANSFER_DECIMALS
                );
            } else {
                Swap memory swap = item.swap;

                address executorAddress = AddressUtils.bytes32ToAddress(swap.executor);

                require(executorAddress != address(0), "Invalid executor");
                require(PeerLibrary.hasStableCoin(self, swap.stableToken), "Invalid stable coin");
                require(swap.command.length > 0, "Invalid command");

                IERC20Metadata stableToken = IERC20Metadata(AddressUtils.bytes32ToAddress(swap.stableToken));
                uint256 stableCoinBalance = stableToken.balanceOf(address(this));
                bytes memory swapResult;

                if (inTransferAddress == address(0)) {
                    require(msg.value >= item.amount, "Invalid amount");
                    swapResult = Address.functionCallWithValue(executorAddress, swap.command, item.amount);
                } else {
                    IERC20Metadata inputToken = IERC20Metadata(inTransferAddress);
                    SafeERC20.safeTransferFrom(inputToken, msg.sender, address(this), item.amount);

                    inputToken.approve(executorAddress, item.amount);
                    swapResult = Address.functionCall(executorAddress, swap.command);
                    inputToken.approve(executorAddress, 0);
                }

                uint256 swapNetUSD = stableToken.balanceOf(address(this)) - stableCoinBalance;
                netUSDX = netUSDX + MathUtils.changeDecimals(
                    swapNetUSD,
                    stableToken.decimals(),
                    CHAIN_TRANSFER_DECIMALS
                );
            }
        }

        require(netUSDX > 0, "Invalid amount");
        require(netUSDX >= _params.inbound.amountUSDX, "Invalid amount");

        /*********** [Estimate gas usage] ***********/
        EstimateOutboundTransfer[]memory estimateOutbounds = new EstimateOutboundTransfer[](_params.outbounds.length);
        for (uint256 i = 0; i < _params.outbounds.length; i++) {
            BridgeOutboundTransfer memory bridgeOutbound = _params.outbounds[i];
            EstimateOutboundTransferItem[]memory items = new EstimateOutboundTransferItem[](bridgeOutbound.items.length);

            for (uint256 j = 0; j < bridgeOutbound.items.length; j++) {
                items[j] = EstimateOutboundTransferItem({
                    asset: bridgeOutbound.items[i].asset,
                    executionGas: bridgeOutbound.items[i].executionGas,
                    commandLength: uint16(bridgeOutbound.items[i].swap.command.length),
                    metadataLength: uint16(bridgeOutbound.items[i].swap.metadata.length)
                });
            }

            estimateOutbounds[i] = EstimateOutboundTransfer({
                chainId: bridgeOutbound.chainId,
                metadataLength: uint16(bridgeOutbound.metadata.length),
                items: items
            });
        }

        EstimateRequest memory estimateRequest = EstimateRequest({
            inbound: EstimateInboundTransfer({
            amountUSDX: netUSDX
        }),
            outbounds: estimateOutbounds
        });

        EstimateResponse memory estimateResult = estimate(estimateRequest);

        netUSDX = netUSDX - (estimateResult.transferGasUSDX + estimateResult.commissionUSDX);
        netUSDX = netUSDX + (estimateResult.transferGasDiscountUSDX + estimateResult.commissionDiscountUSDX);

        EstimateOutboundDetails[] memory outboundMetaData = estimateResult.outboundDetails;

        /*********** [Send transfers] ***********/
        string[] memory messageIds = new string[](_params.outbounds.length);

        for (uint256 i = 0; i < _params.outbounds.length; i++) {
            BridgeOutboundTransfer memory bridgeOutbound = _params.outbounds[i];
            ChainTransferItem[]memory items = new ChainTransferItem[](bridgeOutbound.items.length);
            uint256 chainTotalAmountUSDX = 0;

            for (uint256 j = 0; j < bridgeOutbound.items.length; j++) {
                BridgeOutboundTransferItem memory bridgeItem = bridgeOutbound.items[j];

                uint256 amountUSDX = MathUtils.bps(netUSDX, bridgeItem.distributionBP);
                chainTotalAmountUSDX = chainTotalAmountUSDX + amountUSDX;

                items[i] = ChainTransferItem({
                    amountUSDX: uint64(amountUSDX),
                    asset: bridgeItem.asset,
                    recipient: bridgeItem.recipient,
                    executor: bridgeItem.swap.executor,
                    stableToken: bridgeItem.swap.stableToken,
                    command: bridgeItem.swap.command,
                    metadata: bridgeItem.swap.metadata
                });
            }

            ChainTransfer memory chainTransfer = ChainTransfer({
                items: items,
                metadata: bridgeOutbound.metadata
            });

            messageIds[i] = send(
                bridgeOutbound.chainId,
                MessageTypeChainTransfer,
                ChainTransferLibrary.encode(chainTransfer),
                uint128(outboundMetaData[i].executionGas),
                uint128(outboundMetaData[i].transferGas)
            );

            emit SendChainTransferEvent(
                msg.sender,
                bridgeOutbound.chainId,
                messageIds[i],
                chainTotalAmountUSDX
            );
        }

        return BridgeResponse({
            messageIds: messageIds
        });
    }

    function validateBridge(
        BridgeRequest memory _params
    ) internal view {
        require(_params.inbound.amountUSDX > 0, "Invalid amount");
        require(ChainId.unwrap(gringottsAgents[chainId].chainID) > 0, "Chain not found");

        /************ [Validate inbound transfer] ***********/
        Peer memory self = gringottsAgents[chainId];

        for (uint256 i = 0; i < _params.inbound.items.length; i++) {
            BridgeInboundTransferItem memory item = _params.inbound.items[i];

            require(item.amount > 0, "Invalid amount");
            require(GringottsAddress.unwrap(item.asset) != bytes32(0), "Invalid asset");

            if (PeerLibrary.hasStableCoin(self, item.asset) == false) {
                require(GringottsAddress.unwrap(item.swap.executor) != bytes32(0), "Invalid executor");
                require(item.swap.command.length > 0, "Invalid command");

                require(GringottsAddress.unwrap(item.swap.stableToken) != bytes32(0), "Invalid executor");
                require(PeerLibrary.hasStableCoin(self, item.swap.stableToken), "Invalid stable coin");
            }
        }

        /************ [Validate outbound transfer] ***********/
        uint16 totalDistributionBP = 0;
        uint8 transfers = 0;

        for (uint256 i = 0; i < _params.outbounds.length; i++) {
            BridgeOutboundTransfer memory outbound = _params.outbounds[i];
            require(ChainId.unwrap(outbound.chainId) > 0, "Chain not found");

            Peer memory agent = gringottsAgents[outbound.chainId];
            require(ChainId.unwrap(agent.chainID) > 0, "Chain not found");

            for (uint256 j = 0; j < outbound.items.length; j++) {
                BridgeOutboundTransferItem memory item = outbound.items[j];

                require(item.distributionBP > 0, "Invalid distributionBP");
                require(GringottsAddress.unwrap(item.recipient) != bytes32(0), "Invalid recipient");

                if (PeerLibrary.hasStableCoin(agent, item.asset) == false) {
                    require(GringottsAddress.unwrap(item.swap.executor) != bytes32(0), "Invalid executor");
                    require(item.swap.command.length > 0, "Invalid command");

                    require(GringottsAddress.unwrap(item.swap.stableToken) != bytes32(0), "Invalid executor");
                    require(PeerLibrary.hasStableCoin(agent, item.swap.stableToken), "Invalid stable coin");
                }

                totalDistributionBP = totalDistributionBP + item.distributionBP;
                transfers++;
            }
        }

        require(totalDistributionBP == 10000, "Invalid distributionBP");
        require(transfers <= MAX_TRANSFER, "Too many transfers");
    }

    struct EstimateInboundTransfer {
        uint256 amountUSDX;
    }

    struct EstimateOutboundTransferItem {
        GringottsAddress asset;
        uint256 executionGas;
        uint16 commandLength;
        uint16 metadataLength;
    }

    struct EstimateOutboundTransfer {
        ChainId chainId;
        uint16 metadataLength;
        EstimateOutboundTransferItem[] items;
    }

    struct EstimateRequest {
        EstimateInboundTransfer inbound;
        EstimateOutboundTransfer[] outbounds;
    }

    struct EstimateOutboundDetails {
        ChainId chainId;

        uint256 executionGas;
        uint256 executionGasUSDX;
        uint256 transferGas;
        uint256 transferGasUSDX;
    }

    struct EstimateResponse {
        uint256 commissionUSDX;
        uint256 commissionDiscountUSDX;
        uint256 transferGasUSDX;
        uint256 transferGasDiscountUSDX;

        EstimateOutboundDetails[] outboundDetails;
    }

    function estimate(
        EstimateRequest memory _params
    ) public view returns (EstimateResponse memory) {
        uint256 totalTransferGasPrice = 0;

        EstimateOutboundDetails[] memory outboundDetails = new EstimateOutboundDetails[](_params.outbounds.length);
        uint8 totalTransfers = 0;

        for (uint256 i = 0; i < _params.outbounds.length; i++) {
            EstimateOutboundTransfer memory outbound = _params.outbounds[i];
            require(ChainId.unwrap(outbound.chainId) > 0, "Chain not found");

            Peer memory agent = gringottsAgents[outbound.chainId];
            require(ChainId.unwrap(agent.chainID) > 0, "Chain not found");

            ChainTransferItem[] memory chainTransferItems = new ChainTransferItem[](outbound.items.length);

            uint256 chainExecutionGasPrice = agent.baseGasEstimate;

            for (uint256 j = 0; j < outbound.items.length; j++) {
                EstimateOutboundTransferItem memory item = outbound.items[j];

                chainTransferItems[i] = ChainTransferItem({
                    amountUSDX: 0,
                    asset: GringottsAddress.wrap(bytes32(0)),
                    recipient: GringottsAddress.wrap(bytes32(0)),
                    executor: GringottsAddress.wrap(bytes32(0)),
                    stableToken: GringottsAddress.wrap(bytes32(0)),
                    command: new bytes(item.commandLength),
                    metadata: new bytes(item.metadataLength)
                });

                chainExecutionGasPrice = chainExecutionGasPrice + item.executionGas;
                totalTransfers++;
            }

            ChainTransfer memory chainTransfer = ChainTransfer({
                items: chainTransferItems,
                metadata: new bytes(outbound.metadataLength)
            });

            uint256 transferGasPrice = quote(
                outbound.chainId,
                MessageTypeChainTransfer,
                ChainTransferLibrary.encode(chainTransfer),
                uint128(chainExecutionGasPrice)
            );

            outboundDetails[i] = EstimateOutboundDetails({
                chainId: outbound.chainId,
                executionGas: chainExecutionGasPrice,
                executionGasUSDX: getNativePriceUSD(chainExecutionGasPrice, CHAIN_TRANSFER_DECIMALS),
                transferGas: transferGasPrice,
                transferGasUSDX: getNativePriceUSD(transferGasPrice, CHAIN_TRANSFER_DECIMALS)
            });

            totalTransferGasPrice = totalTransferGasPrice + transferGasPrice;
        }

        require(totalTransfers <= MAX_TRANSFER, "Too many transfers");

        uint256 commissionUSDX = MathUtils.microBPS(_params.inbound.amountUSDX, config.commissionMicroBPS);
        uint256 transferGasUSDX = getNativePriceUSD(totalTransferGasPrice, CHAIN_TRANSFER_DECIMALS);

        uint256 commissionDiscountUSDX = MathUtils.bps(commissionUSDX, config.commissionDiscountBPS);
        uint256 transferGasDiscountUSDX = MathUtils.bps(transferGasUSDX, config.gasDiscountBPS);

        require(_params.inbound.amountUSDX >= (commissionUSDX + transferGasUSDX) - (commissionDiscountUSDX + transferGasDiscountUSDX), "Invalid amount");

        return EstimateResponse({
            commissionUSDX: commissionUSDX,
            commissionDiscountUSDX: commissionDiscountUSDX,
            transferGasUSDX: transferGasUSDX,
            transferGasDiscountUSDX: transferGasDiscountUSDX,
            outboundDetails: outboundDetails
        });
    }

    function _onReceive(
        string memory _guid,
        ChainId _chainId,
        uint8 _header,
        bytes memory _message
    ) internal override {
        _requireNotPaused();

        if (_header == MessageTypeChainTransfer) {
            ChainTransfer memory transfer = ChainTransferLibrary.decode(_message);
            Peer memory self = gringottsAgents[chainId];

            for (uint256 i = 0; i < transfer.items.length; i++) {
                ChainTransferItem memory item = transfer.items[i];

                address recipient = AddressUtils.bytes32ToAddress(item.recipient);

                if (PeerLibrary.hasStableCoin(self, item.asset)) {
                    address stableTokenAddress = AddressUtils.bytes32ToAddress(item.asset);
                    IERC20Metadata stableToken = IERC20Metadata(stableTokenAddress);

                    uint256 amountUSD = MathUtils.changeDecimals(
                        item.amountUSDX,
                        CHAIN_TRANSFER_DECIMALS,
                        stableToken.decimals()
                    );

                    SafeERC20.safeTransfer(stableToken, recipient, amountUSD);
                } else {
                    address stableTokenAddress = AddressUtils.bytes32ToAddress(item.stableToken);
                    address executorAddress = AddressUtils.bytes32ToAddress(item.executor);

                    IERC20Metadata stableToken = IERC20Metadata(stableTokenAddress);

                    uint256 amountUSD = MathUtils.changeDecimals(
                        item.amountUSDX,
                        CHAIN_TRANSFER_DECIMALS,
                        stableToken.decimals()
                    );

                    stableToken.approve(executorAddress, amountUSD);
                    // Can't use openzeppelin's Address library because "Try can only be used with external function calls and contract creation calls"
                    (bool success,) = executorAddress.call(item.command);
                    stableToken.approve(executorAddress, 0);

                    if (!success) {
                        SafeERC20.safeTransfer(stableToken, recipient, amountUSD);
                    }
                }

                emit ReceiveChainTransferEvent(_chainId, _guid, AddressUtils.bytes32ToAddress(item.asset), recipient, item.amountUSDX);
            }
        }
    }

    function updateAgent(
        ChainId _chainID,
        Peer calldata _agent
    ) external onlyOwner {
        require(ChainId.unwrap(_chainID) > 0, "Invalid chain id");

        gringottsAgents[_chainID] = _agent;
        setChainMapping(_chainID, _agent.lzEID);

        if (ChainId.unwrap(chainId) != ChainId.unwrap(_chainID)) {
            setPeer(_agent.lzEID, GringottsAddress.unwrap(_agent.endpoint));
        }
    }

    function setConfig(Config calldata _config) external onlyOwner {
        require(_config.commissionMicroBPS > 0, "Invalid commission");
        config = _config;
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }
}
