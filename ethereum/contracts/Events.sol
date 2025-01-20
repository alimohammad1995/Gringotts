// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {ChainId} from "./Types.sol";

event SendChainTransferEvent (
    uint64 indexed id,
    address indexed sender,
    ChainId chainId,
    bytes32 indexed messageId,
    uint256 amountUSDX
);

event ReceiveChainTransferEvent (
    ChainId chainId,
    bytes32 indexed messageId,
    address asset,
    address indexed recipient,
    uint256 amountUSDX
);