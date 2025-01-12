package provider

import (
	"encoding/base64"
	"github.com/go-resty/resty/v2"
	"github.com/holiman/uint256"
	"gringotts/utils"
	"strconv"
)

const MaxAccounts = 8

type Jupiter struct {
}

// https://quote-api.jup.ag/v6/quote?inputMint=So11111111111111111111111111111111111111112&outputMint=EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v&amount=1000000000&slippageBps=100&restrictIntermediateTokens=true

func (o *Jupiter) GetSwap(params *SwapParams) (*Swap, error) {
	quoteRes := o.quote(
		params.FromToken.Address,
		params.ToToken.Address,
		params.Amount,
		params.SlippageBPS,
	)

	swapRes := o.swap(quoteRes, params.Recipient)

	swapInstruction := swapRes["swapInstruction"].(map[string]interface{})

	swapAccounts := swapInstruction["accounts"].([]interface{})
	swapCommand, _ := base64.StdEncoding.DecodeString(swapInstruction["data"].(string))

	outAmount, _ := strconv.ParseFloat(quoteRes["outAmount"].(string), 64)
	outAmountMin, _ := strconv.ParseFloat(quoteRes["otherAmountThreshold"].(string), 64)

	alts := make([]string, 0, len(swapRes["addressLookupTableAddresses"].([]interface{})))
	for _, alt := range swapRes["addressLookupTableAddresses"].([]interface{}) {
		alts = append(alts, alt.(string))
	}

	accounts, metadata := o.createMeta(swapAccounts)

	return &Swap{
		ExecutorAddress: JupiterAddress,
		Command:         utils.ToHex(swapCommand),
		Metadata:        utils.ToHex(metadata),
		OutAmount:       uint256.NewInt(uint64(outAmount)),
		MinOutAmount:    uint256.NewInt(uint64(outAmountMin)),
		AddressLookup:   alts,
		Accounts:        accounts,
	}, nil
}

func (o *Jupiter) createMeta(swapAccounts []interface{}) ([]Account, []byte) {
	metadata := []byte{byte(len(swapAccounts))}
	var accounts []Account

	for _, swapAccount := range swapAccounts {
		pubkey := swapAccount.(map[string]interface{})["pubkey"].(string)
		isWriteable := swapAccount.(map[string]interface{})["isWritable"].(bool)
		accounts = append(accounts, Account{Address: pubkey, IsSigner: false, IsWriteable: isWriteable})
	}

	return accounts, metadata
}

func (o *Jupiter) swap(quoteResponse map[string]interface{}, recipient string) map[string]interface{} {
	body := map[string]interface{}{
		"quoteResponse": quoteResponse,
		"userPublicKey": recipient,
	}

	var result map[string]interface{}
	_, err := resty.New().R().
		SetBody(body).
		SetResult(&result).
		Post("https://quote-api.jup.ag/v6/swap-instructions")

	if err != nil {
		return nil
	}

	return result
}

func (o *Jupiter) quote(from string, to string, amount *uint256.Int, slippageBPS int) map[string]interface{} {
	if len(from) == 0 {
		from = WSOL
	}
	if len(to) == 0 {
		to = WSOL
	}

	queryParams := map[string]string{
		"inputMint":                  from,
		"outputMint":                 to,
		"amount":                     amount.String(),
		"restrictIntermediateTokens": strconv.FormatBool(true),
		"slippageBps":                strconv.FormatInt(int64(slippageBPS), 10),
		"maxAccounts":                strconv.FormatInt(MaxAccounts, 10),
	}

	var result map[string]interface{}
	_, err := resty.New().R().
		SetQueryParams(queryParams).
		SetResult(&result).
		Get("https://quote-api.jup.ag/v6/quote")

	if err != nil {
		return nil
	}

	return result
}

const WSOL = "So11111111111111111111111111111111111111112"
const JupiterAddress = "JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"
