// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {EndpointV2Mock} from "@layerzerolabs/test-devtools-evm-hardhat/contracts/mocks/EndpointV2Mock.sol";

contract MockLayerZero is EndpointV2Mock {
    constructor(uint32 _eid) EndpointV2Mock(_eid) {}
    }
