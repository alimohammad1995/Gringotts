// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {GringottsAddress} from "./Types.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";

    struct ChainTransferItem {
        uint64 amountUSDX;

        GringottsAddress asset;
        GringottsAddress recipient;

        GringottsAddress executor;
        GringottsAddress stableToken;
        bytes command;
        bytes metadata;
    }

    struct ChainTransfer {
        ChainTransferItem[] items;
    }

library ChainTransferLibrary {
    function decode(bytes memory data) internal pure returns (ChainTransfer memory chainTransfer) {
        uint256 offset = 0;

        uint8 numItems = BytesLib.toUint8(data, offset);
        offset += 1;

        chainTransfer.items = new ChainTransferItem[](numItems);

        for (uint8 i = 0; i < numItems; i++) {
            // Read amountUSDX (uint64)
            uint64 amountUSDX = BytesLib.toUint64(data, offset);
            offset += 8;

            // Read asset, recipient, executor, stableToken (bytes32 each)
            bytes32 asset = BytesLib.toBytes32(data, offset);
            offset += 32;

            bytes32 recipient = BytesLib.toBytes32(data, offset);
            offset += 32;

            bytes32 executor = BytesLib.toBytes32(data, offset);
            offset += 32;

            bytes32 stableToken = BytesLib.toBytes32(data, offset);
            offset += 32;

            // Read commandLength (uint16)
            uint16 commandLength = BytesLib.toUint16(data, offset);
            offset += 2;

            // Read command (bytes)
            bytes memory command = BytesLib.slice(data, offset, commandLength);
            offset += commandLength;

            // Read metadataLength (uint16)
            uint16 metadataLength = BytesLib.toUint16(data, offset);
            offset += 2;

            // Read metadata (bytes)
            bytes memory metadata = BytesLib.slice(data, offset, metadataLength);
            offset += metadataLength;

            // Construct ChainTransferItem
            chainTransfer.items[i] = ChainTransferItem({
                amountUSDX: amountUSDX,
                asset: GringottsAddress.wrap(asset),
                recipient: GringottsAddress.wrap(recipient),
                executor: GringottsAddress.wrap(executor),
                stableToken: GringottsAddress.wrap(stableToken),
                command: command,
                metadata: metadata
            });
        }
    }

    function encode(ChainTransfer memory chainTransfer) internal pure returns (bytes memory) {
        bytes memory encoded = abi.encodePacked(uint8(chainTransfer.items.length));

        for (uint256 i = 0; i < chainTransfer.items.length; i++) {
            ChainTransferItem memory item = chainTransfer.items[i];

            encoded = abi.encodePacked(
                encoded,
                item.amountUSDX,              // uint64 (8 bytes)
                item.asset,                   // bytes32 (32 bytes)
                item.recipient,               // bytes32 (32 bytes)
                item.executor,                // bytes32 (32 bytes)
                item.stableToken,             // bytes32 (32 bytes)
                uint16(item.command.length),  // uint16 (2 bytes)
                item.command,                 // bytes (command bytes)
                uint16(item.metadata.length), // uint16 (2 bytes)
                item.metadata                 // bytes (metadata bytes)
            );
        }

        return encoded;
    }
}