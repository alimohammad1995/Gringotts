// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

type ChainId is uint8;
type GringottsAddress is bytes32;

struct Config {
    uint32 commissionMicroBPS;
    uint32 commissionDiscountBPS;
    uint32 gasDiscountBPS;

    GringottsAddress[] stableCoins;
}

library ConfigLibrary {
    function hasStableCoin(Config memory config, GringottsAddress stableCoin) internal pure returns (bool) {
        for (uint256 i = 0; i < config.stableCoins.length; i++) {
            if (GringottsAddress.unwrap(config.stableCoins[i]) == GringottsAddress.unwrap(stableCoin)) {
                return true;
            }
        }

        return false;
    }
}