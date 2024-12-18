// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "solidity-bytes-utils/contracts/BytesLib.sol";

uint8 constant MessageTypeChainTransfer = 1;

    struct Message {
        uint8 header;
        bytes payload;
    }

library MessageLibrary {
    function create(
        uint8 _header,
        bytes memory _payload
    ) internal pure returns (Message memory) {
        return Message({header: _header, payload: _payload});
    }

    function encode(
        Message memory _message
    ) internal pure returns (bytes memory) {
        return abi.encodePacked(_message.header, _message.payload);
    }

    function decode(
        bytes memory _message
    ) internal pure returns (Message memory) {
        uint8 header = BytesLib.toUint8(_message, 0);
        bytes memory payload = BytesLib.slice(_message, 1, _message.length - 1);
        return Message({header: header, payload: payload});
    }
}
