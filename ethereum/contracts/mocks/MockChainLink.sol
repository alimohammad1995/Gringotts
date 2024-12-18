// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract MockChainLink is AggregatorV3Interface {
    function decimals() external pure returns (uint8) {
        return 10;
    }

    function description() external pure returns (string memory) {
        return "chainlink";
    }

    function version() external pure returns (uint256) {
        return 1;
    }

    function getRoundData(
        uint80
    )
    external
    view
    returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    )
    {
        return (1, 1000, 1, block.timestamp, 1);
    }

    function latestRoundData()
    external
    view
    returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    )
    {
        return (1, 1000, 1, block.timestamp, 1);
    }

    fallback() external payable {}

    receive() external payable {}
}
