package provider

import (
	"errors"
	"fmt"
	"github.com/holiman/uint256"
	"gringotts/config"
	"gringotts/utils"
	"math"
	"strconv"

	"github.com/go-resty/resty/v2"

	"gringotts/models"
)

type OpenOcean struct {
}

// https://open-api.openocean.finance/v3/eth/swap_quote?gasPrice=3&inTokenAddress=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48&outTokenAddress=0xEeeeeEeeeEeEeeEeEeeeeEEEeeeeEeeeeeeeEEeE&amount=500&slippage=1&account=0x929B44e589AC4dD99c0282614e9a844Ea9483aaa&sender=0x0D595AE2666a2c5Ae6b99cce4DD428a9Cf20B2c9

func (o *OpenOcean) GetSwap(params *SwapParams) (*Swap, error) {
	if params.FromToken == nil || params.ToToken == nil {
		return nil, errors.New("invalid token")
	}

	res := o.fetchOpenOceanSwap(params)

	if res == nil || len(res) == 0 {
		return nil, errors.New("open ocean swap not found")
	}

	data := res["data"].(map[string]interface{})

	outAmount, _ := strconv.ParseFloat(data["outAmount"].(string), 64)
	minOutAmount, _ := strconv.ParseFloat(data["minOutAmount"].(string), 64)

	return &Swap{
		ExecutorAddress: data["to"].(string),
		Command:         data["data"].(string),
		OutAmount:       uint256.NewInt(uint64(outAmount)),
		MinOutAmount:    uint256.NewInt(uint64(minOutAmount)),
	}, nil
}

func (o *OpenOcean) fetchOpenOceanSwap(params *SwapParams) map[string]interface{} {
	from := params.FromToken.Address
	to := params.ToToken.Address

	if len(from) == 0 {
		from = OpenOceanNativeMapping[params.Chain]
	}
	if len(to) == 0 {
		to = OpenOceanNativeMapping[params.Chain]
	}

	amount := utils.RoundDown(params.Amount.Float64()/math.Pow10(params.FromToken.Decimals), config.PriceDecimals)

	queryParams := map[string]string{
		"inTokenAddress":  from,
		"outTokenAddress": to,
		"amount":          strconv.FormatFloat(amount, 'f', -1, 64),
		"gasPrice":        "15",
		"slippage":        strconv.FormatFloat(float64(params.SlippageBPS)/100, 'f', -1, 64),
		"sender":          params.Chain.GetContract(),
		"account":         params.Recipient,
		"enabledDexIds":   OpenOceanDexMapping[params.Chain],
	}

	var result map[string]interface{}
	_, err := resty.New().R().
		SetQueryParams(queryParams).
		SetResult(&result).
		Get(fmt.Sprintf("https://open-api.openocean.finance/v3/%s/swap_quote", OpenOceanBlockChainMapping[params.Chain]))

	if err != nil || int(result["code"].(float64)) != 200 {
		return nil
	}

	return result
}

var OpenOceanBlockChainMapping = map[models.Blockchain]string{
	models.Ethereum: "eth",
	models.Binance:  "bsc",
	models.Polygon:  "polygon",
	models.BASE:     "base",
	models.Fantom:   "fantom",
	models.Avax:     "avax",
	models.Arbitrum: "arbitrum",
	models.Optimism: "optimism",
	models.Aurora:   "aurora",
}

var OpenOceanNativeMapping = map[models.Blockchain]string{
	models.Ethereum: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	models.Binance:  "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	models.Polygon:  "0x0000000000000000000000000000000000001010",
	models.BASE:     "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	models.Fantom:   "0x0000000000000000000000000000000000000000",
	models.Avax:     "0x0000000000000000000000000000000000000000",
	models.Arbitrum: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	models.Optimism: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	models.Aurora:   "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
}

var OpenOceanDexMapping = map[models.Blockchain]string{
	models.Ethereum: "1,3,4",
	models.Binance:  "0,1",
	models.Polygon:  "1,2,15,6,37",
	models.BASE:     "1,2,3",
	models.Fantom:   "1,2,3",
	models.Avax:     "1,40,41",
	models.Arbitrum: "1,3,4",
	models.Optimism: "1,2",
	models.Aurora:   "1,9,10",
}
