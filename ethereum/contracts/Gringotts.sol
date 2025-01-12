// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Blockable} from "./utils/Blockable.sol";
import {MessageTypeChainTransfer, MessageTypeChainRegisterOrder, MessageTypeChainCompletion} from "./Message.sol";
import {ChainTransfer, ChainTransferItem, ChainTransferLibrary, EVMTransferItem, EVMTransferLibrary} from "./ChainTransfer.sol";
import {Config, ConfigLibrary, GringottsAddress, ChainId} from "./Types.sol";
import {Peer} from "./Peer.sol";
import {Market} from "./Market.sol";
import {Vault} from "./Vault.sol";
import {AddressUtils, MathUtils, TextUtils} from "./utils/Utils.sol";
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
    uint private constant MULTI_SEND_GAS_INCREASE = 50;

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

    struct BridgeInboundTransferItem {
        GringottsAddress asset;
        uint256 amount;

        // Swap needed
        GringottsAddress executor;
        bytes command;
        GringottsAddress stableToken;
    }

    struct BridgeInboundTransfer {
        uint256 amountUSDX;
        BridgeInboundTransferItem[] items;
    }

    struct BridgeOutboundTransferItem {
        uint16 distributionBP;
    }

    struct BridgeOutboundTransfer {
        ChainId chainId;
        uint256 executionGas;
        bytes message;
        BridgeOutboundTransferItem[] items;
    }

    struct BridgeRequest {
        BridgeInboundTransfer inbound;
        BridgeOutboundTransfer[] outbounds;
    }

    struct BridgeResponse {
        bytes32[] messageIds;
    }

    function bridge(
        BridgeRequest memory _params
    ) external onlyNotBlocked whenNotPaused payable returns (BridgeResponse memory) {
        validateBridge(_params);

        uint256 netUSDX = 0;

        /*********** [Inbound transaction] ***********/
        for (uint256 i = 0; i < _params.inbound.items.length; i++) {
            BridgeInboundTransferItem memory item = _params.inbound.items[i];
            address inTransferAddress = AddressUtils.bytes32ToAddress(item.asset);

            if (ConfigLibrary.hasStableCoin(config, item.asset)) {
                IERC20Metadata stableToken = IERC20Metadata(inTransferAddress);
                SafeERC20.safeTransferFrom(stableToken, msg.sender, address(this), item.amount);

                netUSDX = netUSDX + MathUtils.changeDecimals(
                    item.amount,
                    stableToken.decimals(),
                    CHAIN_TRANSFER_DECIMALS
                );
            } else {
                address executorAddress = AddressUtils.bytes32ToAddress(item.executor);

                IERC20Metadata stableToken = IERC20Metadata(AddressUtils.bytes32ToAddress(item.stableToken));
                uint256 stableCoinBalance = stableToken.balanceOf(address(this));
                bytes memory swapResult;

                if (inTransferAddress == address(0)) {
                    require(msg.value >= item.amount, "Invalid amount");
                    swapResult = Address.functionCallWithValue(executorAddress, item.command, item.amount);
                } else {
                    IERC20Metadata inputToken = IERC20Metadata(inTransferAddress);
                    SafeERC20.safeTransferFrom(inputToken, msg.sender, address(this), item.amount);

                    inputToken.approve(executorAddress, item.amount);
                    swapResult = Address.functionCall(executorAddress, item.command);
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

            estimateOutbounds[i] = EstimateOutboundTransfer({
                chainId: bridgeOutbound.chainId,
                executionGas: bridgeOutbound.executionGas,
                messageLength: uint16(bridgeOutbound.message.length)
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
        bytes32[] memory messageIds = new bytes32[](_params.outbounds.length);

        for (uint256 i = 0; i < _params.outbounds.length; i++) {
            BridgeOutboundTransfer memory bridgeOutbound = _params.outbounds[i];

            Peer memory agent = gringottsAgents[bridgeOutbound.chainId];

            ChainTransferItem[]memory items = new ChainTransferItem[](bridgeOutbound.items.length);
            uint256 chainTotalAmountUSDX = 0;

            for (uint256 j = 0; j < bridgeOutbound.items.length; j++) {
                BridgeOutboundTransferItem memory bridgeItem = bridgeOutbound.items[j];

                uint256 amountUSDX = MathUtils.bps(netUSDX, bridgeItem.distributionBP);
                chainTotalAmountUSDX = chainTotalAmountUSDX + amountUSDX;

                items[i] = ChainTransferItem({
                    amountUSDX: uint64(amountUSDX)
                });
            }

            ChainTransfer memory chainTransfer = ChainTransfer({
                items: items,
                message: bridgeOutbound.message
            });

            if (agent.multiSend) {
                messageIds[i] = send(
                    bridgeOutbound.chainId,
                    MessageTypeChainRegisterOrder,
                    ChainTransferLibrary.encode(chainTransfer),
                    uint128(agent.registerGasEstimate),
                    uint128(outboundMetaData[i].transferGas)
                );

                send(
                    bridgeOutbound.chainId,
                    MessageTypeChainTransfer,
                    abi.encodePacked(messageIds[i]),
                    uint128(outboundMetaData[i].executionGas),
                    uint128(outboundMetaData[i].transferGas)
                );

                send(
                    bridgeOutbound.chainId,
                    MessageTypeChainCompletion,
                    abi.encodePacked(messageIds[i]),
                    uint128(agent.completionGasEstimate),
                    uint128(outboundMetaData[i].transferGas)
                );
            } else {
                messageIds[i] = send(
                    bridgeOutbound.chainId,
                    MessageTypeChainTransfer,
                    ChainTransferLibrary.encode(chainTransfer),
                    uint128(outboundMetaData[i].executionGas),
                    uint128(outboundMetaData[i].transferGas)
                );
            }

            emit SendChainTransferEvent(
                msg.sender,
                bridgeOutbound.chainId,
                TextUtils.toBase64(messageIds[i]),
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

        /************ [Validate inbound transfer] ***********/
        for (uint256 i = 0; i < _params.inbound.items.length; i++) {
            BridgeInboundTransferItem memory item = _params.inbound.items[i];

            require(item.amount > 0, "Invalid amount");
            require(GringottsAddress.unwrap(item.asset) != bytes32(0), "Invalid asset");

            if (ConfigLibrary.hasStableCoin(config, item.asset) == false) {
                require(GringottsAddress.unwrap(item.executor) != bytes32(0), "Invalid executor");
                require(item.command.length > 0, "Invalid command");
                require(GringottsAddress.unwrap(item.stableToken) != bytes32(0), "Invalid executor");
                require(ConfigLibrary.hasStableCoin(config, item.stableToken), "Invalid stable coin");
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

            require(outbound.executionGas > 0, "Invalid execution gas");
            require(outbound.message.length > 0, "Invalid message");

            for (uint256 j = 0; j < outbound.items.length; j++) {
                BridgeOutboundTransferItem memory item = outbound.items[j];

                require(item.distributionBP > 0, "Invalid distributionBP");
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

    struct EstimateOutboundTransfer {
        ChainId chainId;
        uint256 executionGas;
        uint16 messageLength;
    }

    struct EstimateRequest {
        EstimateInboundTransfer inbound;
        EstimateOutboundTransfer[] outbounds;
    }

    struct EstimateOutboundDetails {
        ChainId chainId;

        uint256 executionGas;
        uint256 transferGas;
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

        for (uint256 i = 0; i < _params.outbounds.length; i++) {
            EstimateOutboundTransfer memory outbound = _params.outbounds[i];
            require(ChainId.unwrap(outbound.chainId) > 0, "Chain not found");

            Peer memory agent = gringottsAgents[outbound.chainId];
            require(ChainId.unwrap(agent.chainID) > 0, "Chain not found");

            uint256 transferGasPrice = 0;
            uint256 chainExecutionGasPrice = agent.baseGasEstimate + outbound.executionGas;

            if (agent.multiSend) {
                transferGasPrice += quote(
                    outbound.chainId,
                    MessageTypeChainRegisterOrder,
                    new bytes(outbound.messageLength),
                    uint128(agent.registerGasEstimate)
                );

                transferGasPrice += quote(
                    outbound.chainId,
                    MessageTypeChainTransfer,
                    new bytes(32),
                    uint128(chainExecutionGasPrice)
                );

                transferGasPrice += quote(
                    outbound.chainId,
                    MessageTypeChainCompletion,
                    new bytes(32),
                    uint128(agent.completionGasEstimate)
                );

                chainExecutionGasPrice += agent.registerGasEstimate + agent.completionGasEstimate;
            } else {
                transferGasPrice += quote(
                    outbound.chainId,
                    MessageTypeChainTransfer,
                    new bytes(outbound.messageLength),
                    uint128(chainExecutionGasPrice)
                );
            }

            outboundDetails[i] = EstimateOutboundDetails({
                chainId: outbound.chainId,
                executionGas: chainExecutionGasPrice,
                transferGas: transferGasPrice
            });

            totalTransferGasPrice = totalTransferGasPrice + transferGasPrice;
        }

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
    ) internal whenNotPaused override {
        if (_header == MessageTypeChainTransfer) {
            ChainTransfer memory transfer = ChainTransferLibrary.decode(_message);
            EVMTransferItem[] memory items = EVMTransferLibrary.decode(config, transfer.message);

            require(transfer.items.length == items.length, "Invalid items");

            for (uint256 i = 0; i < transfer.items.length; i++) {
                ChainTransferItem memory chainItem = transfer.items[i];
                EVMTransferItem memory evmItem = items[i];

                address recipient = AddressUtils.bytes32ToAddress(evmItem.recipient);

                if (ConfigLibrary.hasStableCoin(config, evmItem.asset)) {
                    address stableTokenAddress = AddressUtils.bytes32ToAddress(evmItem.asset);
                    IERC20Metadata stableToken = IERC20Metadata(stableTokenAddress);

                    uint256 amountUSD = MathUtils.changeDecimals(
                        chainItem.amountUSDX,
                        CHAIN_TRANSFER_DECIMALS,
                        stableToken.decimals()
                    );

                    SafeERC20.safeTransfer(stableToken, recipient, amountUSD);
                } else {
                    address stableTokenAddress = AddressUtils.bytes32ToAddress(evmItem.stableToken);
                    address executorAddress = AddressUtils.bytes32ToAddress(evmItem.executor);

                    IERC20Metadata stableToken = IERC20Metadata(stableTokenAddress);

                    uint256 amountUSD = MathUtils.changeDecimals(
                        chainItem.amountUSDX,
                        CHAIN_TRANSFER_DECIMALS,
                        stableToken.decimals()
                    );

                    stableToken.approve(executorAddress, amountUSD);
                    // Can't use openzeppelin's Address library because "Try can only be used with external function calls and contract creation calls"
                    (bool success,) = executorAddress.call(evmItem.command);
                    stableToken.approve(executorAddress, 0);

                    if (!success) {
                        SafeERC20.safeTransfer(stableToken, recipient, amountUSD);
                    }
                }

                emit ReceiveChainTransferEvent(_chainId, _guid, AddressUtils.bytes32ToAddress(evmItem.asset), recipient, chainItem.amountUSDX);
            }
        }
    }

    event TestMessage (bytes32 messageId);

    function testSend(ChainId _chainID, uint8 _header, uint8 _header2, bytes memory _message, uint128 _gas, bool multi) external {
        if (multi) {
            bytes32 x = send(
                _chainID,
                _header,
                _message,
                _gas,
                uint128(quote(_chainID, _header, _message, _gas))
            );
            emit TestMessage(x);

            bytes32 y = send(
                _chainID,
                _header2,
                abi.encodePacked(x),
                _gas,
                uint128(quote(_chainID, _header, abi.encodePacked(x), _gas))
            );
            emit TestMessage(y);
        } else {
            bytes32 x = send(
                _chainID,
                _header,
                _message,
                _gas,
                uint128(quote(_chainID, _header, _message, _gas))
            );
            emit TestMessage(x);
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
