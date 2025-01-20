// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import {ChainId, GringottsAddress} from "../Types.sol";

library IdUtils {
    function generateId(ChainId _chainId, uint64 _txCount) internal pure returns (uint64) {
        require(_txCount <= 0xFFFFFFFFFFFFFF, "Transaction count exceeds 7 bytes");

        uint64 chainId = uint64(ChainId.unwrap(_chainId));
        return (chainId << 56) | _txCount;
    }
}

library AddressUtils {
    function bytes32ToAddress(GringottsAddress _inAddress) internal pure returns (address) {
        bytes32 _address = GringottsAddress.unwrap(_inAddress);
        require(uint256(_address) >> 160 == 0, "Invalid bytes32 for address");
        return address(uint160(uint256(_address)));
    }
}

library MathUtils {
    function bps(uint256 _amount, uint32 _bps) internal pure returns (uint256) {
        return microBPS(_amount, _bps * 1000);
    }

    function microBPS(uint256 _amount, uint32 _microBps) internal pure returns (uint256) {
        if (_microBps == 0) {
            return 0;
        }

        return _amount * _microBps / 10_000_000;
    }

    function percentage(uint256 _amount, uint8 _percent) internal pure returns (uint256) {
        return _amount * _percent / 100;
    }

    function changeDecimals(uint256 _value, uint8 _currentDecimals, uint8 _newDecimals) internal pure returns (uint256) {
        return _value * (10 ** _newDecimals) / (10 ** _currentDecimals);
    }
}