// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {ChainId} from "./Types.sol";

    event SendChainTransferItemEvent (
        ChainId chainId,
        string indexed messageId,
        uint256 amountUSDX
    );

    event ReceiveChainTransferItemEvent (
        ChainId chainId,
        string indexed messageId,
        address asset,
        address recipient,
        uint256 amountUSDX
    );