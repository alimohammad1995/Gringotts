// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {BytesLib} from "solidity-bytes-utils/contracts/BytesLib.sol";
import {GringottsAddress} from "./Types.sol";

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
    bytes metadata;
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

            // Read asset, recipient, executor, stableToken (bytes32 each)
            bytes32 asset = BytesLib.toBytes32(data, offset);
            offset += 32;

            bytes32 recipient = BytesLib.toBytes32(data, offset);
            offset += 32;

            uint8 needSwap = BytesLib.toUint8(data, offset);
            offset += 1;

            bytes32 executor;
            bytes32 stableToken;
            bytes memory itemCommand;
            bytes memory itemMetadata;

            if (needSwap > 0) {
                executor = BytesLib.toBytes32(data, offset);
                offset += 32;

                stableToken = BytesLib.toBytes32(data, offset);
                offset += 32;

                // Read itemCommandLength (uint16)
                uint16 itemCommandLength = BytesLib.toUint16(data, offset);
                offset += 2;
                // Read command (bytes)
                itemCommand = BytesLib.slice(data, offset, itemCommandLength);
                offset += itemCommandLength;

                // Read itemMetadataLength (uint16)
                uint16 itemMetadataLength = BytesLib.toUint16(data, offset);
                offset += 2;
                // Read metadata (bytes)
                itemMetadata = BytesLib.slice(data, offset, itemMetadataLength);
                offset += itemMetadataLength;
            }

            items[i] = ChainTransferItem({
                amountUSDX: amountUSDX,
                asset: GringottsAddress.wrap(asset),
                recipient: GringottsAddress.wrap(recipient),
                executor: GringottsAddress.wrap(executor),
                stableToken: GringottsAddress.wrap(stableToken),
                command: itemCommand,
                metadata: itemMetadata
            });
        }

        // Read metadataLength (uint16)
        uint16 metadataLength = BytesLib.toUint16(data, offset);
        offset += 2;
        // Read metadata (bytes)
        bytes memory metadata = BytesLib.slice(data, offset, metadataLength);
        offset += metadataLength;

        return ChainTransfer({
            metadata: metadata,
            items: items
        });
    }

    function encode(ChainTransfer memory chainTransfer) internal pure returns (bytes memory) {
        bytes memory encoded = abi.encodePacked(uint8(chainTransfer.items.length));

        for (uint256 i = 0; i < chainTransfer.items.length; i++) {
            ChainTransferItem memory item = chainTransfer.items[i];

            uint8 needSwap = item.command.length > 0 ? 1 : 0;

            encoded = abi.encodePacked(
                encoded,
                item.amountUSDX,                // uint64 (8 bytes)
                item.asset,                     // bytes32 (32 bytes)
                item.recipient,                 // bytes32 (32 bytes)
                needSwap                        // bool (1 bytes)
            );

            if (needSwap > 0) {
                encoded = abi.encodePacked(
                    encoded,
                    item.executor,                  // bytes32 (32 bytes)
                    item.stableToken,               // bytes32 (32 bytes)
                    uint16(item.command.length),    // uint16 (2 bytes)
                    item.command,                   // bytes
                    uint16(item.metadata.length),   // uint16 (2 bytes)
                    item.metadata                   // bytes
                );
            }
        }

        encoded = abi.encodePacked(
            encoded,
            uint16(chainTransfer.metadata.length),  // uint16 (2 bytes)
            chainTransfer.metadata                  // bytes
        );

        return encoded;
    }
}