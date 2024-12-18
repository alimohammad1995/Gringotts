// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/access/Ownable.sol";

abstract contract Blockable is Ownable {
    event BlockEvent(address indexed account, bool blocked);

    mapping(address => bool) private _blockedAccounts;

    modifier onlyNotBlocked() {
        require(_blockedAccounts[msg.sender] == false, "Account is blocked");
        _;
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
