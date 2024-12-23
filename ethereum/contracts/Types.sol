// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

type ChainId is uint8;
type GringottsAddress is bytes32;

struct Config {
    uint32 commissionMicroBPS;
    uint32 commissionDiscountBPS;
    uint32 gasDiscountBPS;
}