package provider

import (
	"errors"
	"github.com/holiman/uint256"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-resty/resty/v2"

	"gringotts/models"
	"gringotts/utils"
)

type SunSwap struct {
}

// https://rot.endjgfsv.link/swap/router?fromToken=T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb&toToken=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t&amountIn=1000000

func (o *SunSwap) GetSwap(params *models.SwapParams) (*models.Swap, error) {
	res := o.fetchSunSwap(params.FromToken.Address, params.ToToken.Address, params.Amount)

	suggestedRoute := res["data"].([]interface{})
	if res == nil || len(suggestedRoute) == 0 {
		return nil, errors.New("failed to fetch quote")
	}

	contractAbi, err := abi.JSON(strings.NewReader(SunSwapABI))
	if err != nil {
		return nil, errors.New("failed to fetch quote")
	}

	bestRoute := suggestedRoute[len(suggestedRoute)-1].(map[string]interface{})

	outAmountRaw, _ := strconv.ParseFloat(bestRoute["amountOut"].(string), 64)
	outAmount := uint256.NewInt(uint64(outAmountRaw * math.Pow10(params.ToToken.Decimals)))
	minAmountOut := utils.ApplyReverseBPS(outAmount, params.SlippageBPS)

	path := make([]common.Address, 0)
	pathFee := make([]*big.Int, 0)
	for i, tokenAddress := range bestRoute["tokens"].([]interface{}) {
		hexAddress, _ := utils.TronAddressToHex(tokenAddress.(string))
		path = append(path, common.HexToAddress(hexAddress))
		currentFee, _ := strconv.Atoi(bestRoute["poolFees"].([]interface{})[i].(string))
		pathFee = append(pathFee, big.NewInt(int64(currentFee)))
	}

	routePoolVersion := bestRoute["poolVersions"].([]interface{})
	poolVersion := []string{routePoolVersion[0].(string)}
	versionLength := []*big.Int{big.NewInt(2)}

	for _, version := range routePoolVersion[1:] {
		if poolVersion[len(poolVersion)-1] == version {
			versionLength[len(versionLength)-1].Add(versionLength[len(versionLength)-1], big.NewInt(1))
		} else {
			poolVersion = append(poolVersion, version.(string))
			versionLength = append(versionLength, big.NewInt(1))
		}
	}

	hexRecipient, err := utils.TronAddressToHex(params.Recipient)
	if err != nil {
		return nil, err
	}

	swapData := struct {
		AmountIn     *big.Int
		AmountOutMin *big.Int
		To           common.Address
		Deadline     *big.Int
	}{
		AmountIn:     params.Amount.ToBig(),
		AmountOutMin: minAmountOut.ToBig(),
		To:           common.HexToAddress(hexRecipient),
		Deadline:     uint256.NewInt(uint64(time.Now().Unix() + 500*60)).ToBig(),
	}

	data, err := contractAbi.Pack("swapExactInput", path, poolVersion, versionLength, pathFee, swapData)
	if err != nil {
		return nil, err
	}

	contractAddress, _ := utils.TronAddressToHex(SunSwapAddress)

	return &models.Swap{
		Executor:     contractAddress,
		Command:      utils.ToHex(data),
		OutAmount:    outAmount,
		MinOutAmount: minAmountOut,
	}, nil
}

func (o *SunSwap) fetchSunSwap(from string, to string, amount *uint256.Int) map[string]interface{} {
	if len(from) == 0 {
		from = TRX
	}
	if len(to) == 0 {
		to = TRX
	}

	queryParams := map[string]string{
		"fromToken": from,
		"toToken":   to,
		"amountIn":  amount.String(),
		"typeList":  "PSM,CURVE,CURVE_COMBINATION,WTRX,SUNSWAP_V1,SUNSWAP_V2,SUNSWAP_V3",
	}

	var result map[string]interface{}
	_, err := resty.New().R().
		SetQueryParams(queryParams).
		SetResult(&result).
		Get("https://rot.endjgfsv.link/swap/router")

	if err != nil || result["message"] != "SUCCESS" || int(result["code"].(float64)) != 0 {
		return nil
	}

	return result
}

const TRX = "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb" // 0x0000000000000000000000000000000
const SunSwapABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_v2Router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_v1Foctroy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_psmUsdd\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_v3Router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_wtrx\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"AddPool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"ChangePool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amountsOut\",\"type\":\"uint256[]\"}],\"name\":\"SwapExactETHForTokens\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amountsOut\",\"type\":\"uint256[]\"}],\"name\":\"SwapExactTokensForTokens\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"originOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"TransferAdminship\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"originOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"TransferOwnership\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"WTRX\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"poolVersion\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"addPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"poolVersion\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"gemJoin\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"addPsmPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"poolVersion\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"addUsdcPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"changePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"poolVersion\",\"type\":\"string\"}],\"name\":\"isPsmPool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"poolVersion\",\"type\":\"string\"}],\"name\":\"isUsdcPool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"psmUsdd\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"retrieve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"poolVersion\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"versionLen\",\"type\":\"uint256[]\"},{\"internalType\":\"uint24[]\",\"name\":\"fees\",\"type\":\"uint24[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"internalType\":\"struct SmartExchangeRouter.SwapData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"swapExactInput\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amountsOut\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminship\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"unwrapWTRX\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"v1Factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"v2Router\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"v3Router\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]\n"
const SunSwapAddress = "TJ4NNy8xZEqsowCBhLvZ45LCqPdGjkET5j"
