// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/access/Ownable.sol";

abstract contract Blockable is Ownable {
    event BlockEvent(address indexed account, bool blocked);

    error BlockedSender();

    mapping(address => bool) private _blockedAccounts;

    modifier onlyNotBlocked() {
        _requireBlocked();
        _;
    }

    function _requireBlocked() internal view virtual {
        if (_blockedAccounts[_msgSender()]) {
            revert BlockedSender();
        }
    }

    function blockAccount(address account) external onlyOwner {
        _blockedAccounts[account] = true;
        emit BlockEvent(account, true);
    }

    function unblockAccount(address account) external onlyOwner {
        delete _blockedAccounts[account];
        emit BlockEvent(account, false);
    }
}
