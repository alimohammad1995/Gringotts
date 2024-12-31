// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

abstract contract Vault is Ownable {
    function withdrawNative(uint256 _amount) external onlyOwner {
        require(address(this).balance > 0, "No native tokens available");

        if (_amount == 0) {
            payable(this.owner()).transfer(address(this).balance);
        } else {
            payable(this.owner()).transfer(_amount);
        }
    }

    function withdrawERC20(address _token, uint256 _amount) external onlyOwner {
        IERC20 token = IERC20(_token);
        uint256 balance = token.balanceOf(address(this));
        require(balance > 0, "No ERC20 tokens available");

        if (_amount == 0) {
            require(token.transfer(this.owner(), balance), "ERC20 token withdraw failed");
        } else {
            require(token.transfer(this.owner(), _amount), "ERC20 token withdraw failed");
        }
    }

    function balanceERC20(address _token) external view returns (uint256) {
        IERC20 token = IERC20(_token);
        return token.balanceOf(address(this));
    }

    fallback() external payable {}

    receive() external payable {}
}
