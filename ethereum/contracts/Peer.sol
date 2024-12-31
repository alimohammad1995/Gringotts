// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {ChainId, GringottsAddress} from "./Types.sol";
import {Gringotts} from "./Gringotts.sol";

struct Peer {
    ChainId chainID;
    GringottsAddress endpoint;
    uint32 lzEID;
    GringottsAddress[] stableCoins;
    uint128 baseGasEstimate;
}

library PeerLibrary {
    function hasStableCoin(Peer memory agent, GringottsAddress stableCoin) internal pure returns (bool) {
        for (uint256 i = 0; i < agent.stableCoins.length; i++) {
            if (GringottsAddress.unwrap(agent.stableCoins[i]) == GringottsAddress.unwrap(stableCoin)) {
                return true;
            }
        }

        return false;
    }
}