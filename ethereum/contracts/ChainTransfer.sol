// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {BytesLib} from "solidity-bytes-utils/contracts/BytesLib.sol";
import {GringottsAddress} from "./Types.sol";

struct ChainTransferItem {
    uint64 amountUSDX;
}

struct ChainTransfer {
    bytes message;
    ChainTransferItem[] items;
}

struct EVMTransferItem {
    GringottsAddress asset;
    GringottsAddress recipient;

    GringottsAddress executor;
    GringottsAddress stableToken;
    bytes command;
}

library ChainTransferLibrary {
    function decode(bytes memory data) internal pure returns (ChainTransfer memory) {
        uint256 offset = 0;

        uint8 numItems = BytesLib.toUint8(data, offset);
        offset += 1;

        ChainTransferItem[] memory items = new ChainTransferItem[](numItems);

        for (uint8 i = 0; i < numItems; i++) {
            // Read amountUSDX (uint64)
            uint64 amountUSDX = BytesLib.toUint64(data, offset);
            offset += 8;

            items[i] = ChainTransferItem({
                amountUSDX: amountUSDX
            });
        }

        // Read messageLength (uint16)
        uint16 messageLength = BytesLib.toUint16(data, offset);
        offset += 2;
        // Read message (bytes)
        bytes memory message = BytesLib.slice(data, offset, messageLength);
        offset += messageLength;

        return ChainTransfer({
            message: message,
            items: items
        });
    }

    function encode(ChainTransfer memory chainTransfer) internal pure returns (bytes memory) {
        bytes memory encoded = abi.encodePacked(uint8(chainTransfer.items.length));

        for (uint256 i = 0; i < chainTransfer.items.length; i++) {
            ChainTransferItem memory item = chainTransfer.items[i];

            encoded = abi.encodePacked(
                encoded,
                item.amountUSDX                // uint64 (8 bytes)
            );
        }

        encoded = abi.encodePacked(
            encoded,
            uint16(chainTransfer.message.length),   // uint16 (2 bytes)
            chainTransfer.message                   // bytes
        );

        return encoded;
    }
}

library EVMTransferLibrary {
    function decode(bytes memory data) internal pure returns (EVMTransferItem[] memory) {
        uint256 offset = 0;

        uint8 numItems = BytesLib.toUint8(data, offset);
        offset += 1;

        EVMTransferItem[] memory items = new EVMTransferItem[](numItems);

        for (uint8 i = 0; i < numItems; i++) {
            // Read asset, recipient, executor, stableToken (bytes32 each)
            bytes32 asset = BytesLib.toBytes32(data, offset);
            offset += 32;

            bytes32 recipient = BytesLib.toBytes32(data, offset);
            offset += 32;

            uint8 needSwap = BytesLib.toUint8(data, offset);
            offset += 1;

            bytes32 executor;
            bytes32 stableToken;
            bytes memory command;

            if (needSwap > 0) {
                executor = BytesLib.toBytes32(data, offset);
                offset += 32;

                stableToken = BytesLib.toBytes32(data, offset);
                offset += 32;

                // Read commandLength (uint16)
                uint16 commandLength = BytesLib.toUint16(data, offset);
                offset += 2;
                // Read command (bytes)
                command = BytesLib.slice(data, offset, commandLength);
                offset += commandLength;
            }

            items[i] = EVMTransferItem({
                asset: GringottsAddress.wrap(asset),
                recipient: GringottsAddress.wrap(recipient),
                executor: GringottsAddress.wrap(executor),
                stableToken: GringottsAddress.wrap(stableToken),
                command: command
            });
        }


        return items;
    }
}