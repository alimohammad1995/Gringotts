// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {ChainId} from "./Types.sol";

    event SendChainTransferEvent (
        address indexed user,
        ChainId chainId,
        string indexed messageId,
        uint256 amountUSDX
    );

    event ReceiveChainTransferEvent (
        ChainId chainId,
        string indexed messageId,
        address asset,
        address recipient,
        uint256 amountUSDX
    );