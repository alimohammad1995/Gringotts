// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import '@uniswap/v2-periphery/contracts/interfaces/IUniswapV2Router02.sol';

import '@uniswap/v3-periphery/contracts/interfaces/ISwapRouter.sol';
import '@uniswap/v3-periphery/contracts/interfaces/IQuoter.sol';
import '@uniswap/v3-periphery/contracts/libraries/TransferHelper.sol';

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract TestTron {
    IUniswapV2Router02 private routerV2;
    ISwapRouter private routerV3;

    address private weth9;
    uint24 public constant feeTier = 3000;

    constructor(address _routerV2, address _routerV3, address _WETH9) {
        routerV2 = IUniswapV2Router02(_routerV2);
        routerV3 = ISwapRouter(_routerV3);
        weth9 = _WETH9;
    }

    function getDecimals(address target) public view returns (uint8) {
        // Encoding the function signature for `decimals()`
        bytes memory data = abi.encodeWithSignature("decimals()");

        // Perform the low-level call
        (bool success, bytes memory returnData) = target.staticcall(data);
        require(success, "Call to decimals() failed");

        // Decode the return data
        uint8 decimals = abi.decode(returnData, (uint8));
        return decimals;
    }

    event Response(bool success, bytes data);

    function execute(address target, bytes memory data, uint256 _value, uint256 _gas) external returns (bool, bytes memory) {
        (bool success, bytes memory data) = target.call{value: _value, gas: _gas}(data);
        emit Response(success, data);
        return (success, data);
    }

    function swapToETHV3(uint256 amountIn, address usdtAddress) external returns (uint256) {
        IERC20(usdtAddress).approve(address(routerV3), amountIn);

        ISwapRouter.ExactInputSingleParams memory params = ISwapRouter.ExactInputSingleParams({
            tokenIn: usdtAddress,
            tokenOut: weth9,
            fee: feeTier,
            recipient: address(this),
            deadline: block.timestamp + 60,
            amountIn: amountIn,
            amountOutMinimum: 0,
            sqrtPriceLimitX96: 0
        });

        uint256 x = routerV3.exactInputSingle(params);
//        IWETH9(weth9).withdraw(x);

        return x;
    }

    function swapFromETHV3(uint256 amountIn, address usdtAddress) external returns (uint256) {
        ISwapRouter.ExactInputSingleParams memory params = ISwapRouter.ExactInputSingleParams({
            tokenIn: weth9,
            tokenOut: usdtAddress,
            fee: feeTier,
            recipient: address(this),
            deadline: block.timestamp + 60,
            amountIn: amountIn,
            amountOutMinimum: 0,
            sqrtPriceLimitX96: 0
        });

        return routerV3.exactInputSingle{value: amountIn}(params);
    }

    function parseRevertReason(bytes memory reason) private pure returns (uint256) {
        if (reason.length != 32) {
            if (reason.length < 68) revert('Unexpected error');
            assembly {
                reason := add(reason, 0x04)
            }
            revert(abi.decode(reason, (string)));
        }
        return abi.decode(reason, (uint256));
    }

    function swapToETHV2(uint256 amountIn, address usdtAddress) external returns (uint[] memory) {
        IERC20(usdtAddress).approve(address(routerV2), amountIn);

        address[] memory path = new address[](2);
        path[0] = usdtAddress;
        path[1] = routerV2.WETH();

        uint256 amountOutMin = 1;
        uint256 deadline = block.timestamp + 60;

        return routerV2.swapExactTokensForETH(
            amountIn,
            amountOutMin,
            path,
            address(this),
            deadline
        );
    }

    function swapFromETHV2(uint256 amountIn, address usdtAddress) external returns (uint[] memory) {
        address[] memory path = new address[](2);
        path[0] = routerV2.WETH();
        path[1] = usdtAddress;

        uint256 amountOutMin = 1;
        uint256 deadline = block.timestamp + 60;

        return routerV2.swapExactETHForTokens{value: amountIn}(
            amountOutMin,
            path,
            address(this),
            deadline
        );
    }

    function estimateV2(address _token, uint256 _amount) public view returns (uint[] memory) {
        address[] memory path = new address[](2);
        path[0] = routerV2.WETH();
        path[1] = _token;

        return routerV2.getAmountsOut(_amount, path);
    }

    fallback() external payable {
    }

    receive() external payable {
    }
}
