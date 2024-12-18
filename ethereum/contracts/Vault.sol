// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

abstract contract Vault is Ownable {
    function withdrawNative() external onlyOwner {
        require(address(this).balance > 0, "No native tokens available");
        payable(this.owner()).transfer(address(this).balance);
    }

    function withdrawERC20(address _token) external onlyOwner {
        IERC20 token = IERC20(_token);
        uint256 balance = token.balanceOf(address(this));
        require(balance > 0, "No ERC20 tokens available");
        require(
            token.transfer(this.owner(), balance),
            "ERC20 token withdraw failed"
        );
    }

    function sendNative(address _to, uint256 _amount) internal {
        require(_to != address(0), "Invalid recipient");
        payable(_to).transfer(_amount);
    }

    function sendERC20(address _token, address _to, uint256 _amount) internal {
        require(_to != address(0), "Invalid recipient");
        IERC20 token = IERC20(_token);
        require(token.transfer(_to, _amount), "ERC20 token send failed");
    }

    function balanceNative() external view returns (uint256) {
        return address(this).balance;
    }

    function balanceERC20(address _token) external view returns (uint256) {
        IERC20 token = IERC20(_token);
        return token.balanceOf(address(this));
    }

    fallback() external payable {}

    receive() external payable {}
}
