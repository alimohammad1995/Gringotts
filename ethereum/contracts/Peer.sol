// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {ChainId, GringottsAddress} from "./Types.sol";
import {Gringotts} from "./Gringotts.sol";

struct Peer {
    ChainId chainID;
    GringottsAddress endpoint;
    uint32 lzEID;

    bool multiSend;

    uint128 baseGasEstimate;
    uint128 registerGasEstimate;
    uint128 completionGasEstimate;
}
