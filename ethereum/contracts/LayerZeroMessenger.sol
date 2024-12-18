// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {ChainId} from "./Types.sol";
import {Message, MessageLibrary} from "./Message.sol";
import {TextUtils} from "./utils/Utils.sol";
import {OApp, Origin, MessagingFee, MessagingReceipt} from "@layerzerolabs/oapp-evm/contracts/oapp/OApp.sol";
import {OptionsBuilder} from "@layerzerolabs/oapp-evm/contracts/oapp/libs/OptionsBuilder.sol";

abstract contract LayerZeroBridge is OApp {
    mapping(ChainId => uint32) private chainEIDMappings;
    mapping(uint32 => ChainId) private eIDChainMappings;

    constructor(address _endpoint, address _owner) OApp(_endpoint, _owner) {
    }

    function send(
        ChainId _chainId,
        uint8 _header,
        bytes memory _message,
        uint128 _receiverGasLimitNative
    ) internal returns (string memory) {
        require(_message.length > 0, "Invalid message");
        require(_receiverGasLimitNative > 0, "Invalid receiver gas limit");
        require(chainEIDMappings[_chainId] > 0, "Invalid chain ID");

        Message memory message = MessageLibrary.create(_header, _message);
        bytes memory messageEncoded = MessageLibrary.encode(message);

        uint256 fee = quote(_chainId, _header, messageEncoded, _receiverGasLimitNative);

        MessagingReceipt memory receipt = _lzSend(
            chainEIDMappings[_chainId],
            messageEncoded,
            _build_options(_receiverGasLimitNative),
            MessagingFee(fee * 2, 0),
            payable(address(this))
        );

        return TextUtils.toBase64(receipt.guid);
    }

    function quote(
        ChainId _chainId,
        uint8 _header,
        bytes memory _message,
        uint128 _receiverGasLimitNative
    ) internal view returns (uint256) {
        require(_message.length > 0, "Invalid message");
        require(_receiverGasLimitNative > 0, "Invalid receiver gas limit");
        require(chainEIDMappings[_chainId] > 0, "Invalid chain ID");

        Message memory message = MessageLibrary.create(_header, _message);
        bytes memory messageEncoded = MessageLibrary.encode(message);

        MessagingFee memory fee = _quote(
            chainEIDMappings[_chainId],
            messageEncoded,
            _build_options(_receiverGasLimitNative),
            false
        );

        return fee.nativeFee;
    }

    function _build_options(
        uint128 _receiverGasLimitNative
    ) private pure returns (bytes memory) {
        return
            OptionsBuilder.addExecutorLzReceiveOption(
            OptionsBuilder.newOptions(),
            _receiverGasLimitNative,
            0
        );
    }

    function _onReceive(
        string memory _guid,
        ChainId _chainId,
        uint8 _header,
        bytes memory _message
    ) internal virtual;

    function _lzReceive(
        Origin calldata _origin,
        bytes32 _guid,
        bytes calldata payload,
        address,
        bytes calldata
    ) internal override {
        ChainId chainId = eIDChainMappings[_origin.srcEid];
        require(ChainId.unwrap(chainId) > 0, "Invalid chain ID");
        Message memory message = MessageLibrary.decode(payload);
        _onReceive(TextUtils.toBase64(_guid), chainId, message.header, message.payload);
    }

    function _payNative(
        uint256 _nativeFee
    ) internal pure override returns (uint256 nativeFee) {
        return _nativeFee;
    }

    function setChainMapping(ChainId _chainId, uint32 _eid) internal onlyOwner {
        require(_eid > 0, "Invalid eid");
        require(ChainId.unwrap(_chainId) > 0, "Invalid chain id");
        chainEIDMappings[_chainId] = _eid;
        eIDChainMappings[_eid] = _chainId;
    }

    function testSend(
        uint32 _dstEid,
        string memory _m
    ) external returns (string memory) {
        bytes memory _message = bytes(_m);

        MessagingFee memory fee = _quote(
            _dstEid,
            _message,
            _build_options(100000),
            false
        );

        MessagingReceipt memory receipt = _lzSend(
            _dstEid,
            _message,
            _build_options(100000),
            MessagingFee(fee.nativeFee * 2, 0),
            payable(address(this))
        );

        return TextUtils.toBase64(receipt.guid);
    }
}
