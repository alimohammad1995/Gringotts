// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {MathUtils} from "./utils/Utils.sol";

import {IPyth} from "@pythnetwork/pyth-sdk-solidity/IPyth.sol";
import {PythStructs} from "@pythnetwork/pyth-sdk-solidity/PythStructs.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

abstract contract Market is Ownable {
    uint256 private constant UPDATE_THRESHOLD = 5 * 60;

    uint8 private networkDecimals;

    address private pythPriceFeed;
    bytes32 private pythPriceFeedId;

    address private chainlinkPriceFeed;

    address private winklinkPriceFeed;

    constructor(uint8 _decimals) {
        networkDecimals = _decimals;
    }

    function getNativePriceUSD(uint256 amount, uint8 desiredDecimals) internal view returns (uint256) {
        uint256 usdValue = 0;

        if (chainlinkPriceFeed != address(0)) {
            usdValue = getPriceChainlink(chainlinkPriceFeed, amount);
        } else if (pythPriceFeed != address(0)) {
            usdValue = getPricePyth(amount);
        } else if (winklinkPriceFeed != address(0)) {
            usdValue = getPriceChainlink(winklinkPriceFeed, amount);
        }

        require(usdValue > 0, "Price not available");
        return MathUtils.changeDecimals(usdValue, networkDecimals, desiredDecimals);
    }

    function getPriceChainlink(address feedAddress, uint256 amount) private view returns (uint256) {
        AggregatorV3Interface feed = AggregatorV3Interface(feedAddress);
        (, int price, , uint256 updatedAt,) = feed.latestRoundData();

        require(updatedAt > block.timestamp - UPDATE_THRESHOLD, "Data feed not updated recently");
        require(price > 0, "Price not available");

        return amount * uint256(price) / (10 ** feed.decimals());
    }

    function getPricePyth(uint256 amount) private view returns (uint256) {
        IPyth priceFeed = IPyth(pythPriceFeed);
        PythStructs.Price memory price = priceFeed.getPriceNoOlderThan(pythPriceFeedId, UPDATE_THRESHOLD);

        uint256 usdValue = amount * (uint64(price.price) + price.conf);
        if (price.expo >= 0) {
            usdValue = uint64(usdValue * 10 ** uint32(price.expo));
        } else {
            usdValue = uint64(usdValue / 10 ** uint32(price.expo));
        }
        return usdValue;
    }

    function setPythPriceFeed(address _pythPriceFeed, bytes32 _pythPriceFeedId) external onlyOwner {
        require(_pythPriceFeed != address(0), "Invalid pyth price feed");
        require(_pythPriceFeedId != bytes32(0), "Invalid pyth price feed id");

        pythPriceFeed = _pythPriceFeed;
        pythPriceFeedId = _pythPriceFeedId;
    }

    function setChainlinkPriceFeed(address _chainlinkPriceFeed) external onlyOwner {
        require(_chainlinkPriceFeed != address(0), "Invalid chainlink price feed");
        chainlinkPriceFeed = _chainlinkPriceFeed;
    }

    function setWinklinkPriceFeed(address _winklinkPriceFeed) external onlyOwner {
        require(_winklinkPriceFeed != address(0), "Invalid winklink price feed");
        winklinkPriceFeed = _winklinkPriceFeed;
    }
}
