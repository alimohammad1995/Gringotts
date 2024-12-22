package models

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
	"os"
)

type Token struct {
	Name         string
	Address      string
	Icon         string
	Decimals     int
	IsStableCoin bool
	IsNative     bool
}

var stableCoins = map[Blockchain][]*Token{}
var tokens = map[Blockchain][]*Token{}

func LoadTokens() error {
	content, err := os.ReadFile(viper.GetString("tokens"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &tokens); err != nil {
		return err
	}

	tokenCount := 0
	for chain, tokenList := range tokens {
		if !chain.IsSupported() {
			return fmt.Errorf("blockchain %s is not supported", chain)
		}

		tokenCount += len(tokenList)
	}
	log.Info(fmt.Sprintf("Loaded %d tokens", tokenCount))

	SetStableCoins()
	tokenCount = 0
	for _, tokenList := range stableCoins {
		tokenCount += len(tokenList)
	}
	log.Info(fmt.Sprintf("Loaded %d stable tokens", tokenCount))

	if len(stableCoins[Solana]) != 1 {
		return fmt.Errorf("solana can only have 1 stable coin")
	}

	return nil
}

func SetStableCoins() {
	if len(stableCoins) != 0 {
		return
	}

	for blockchain, tokensList := range tokens {
		stableTokens := make([]*Token, 0)

		for _, token := range tokensList {
			if token.IsStableCoin {
				stableTokens = append(stableTokens, token)
			}
		}

		stableCoins[blockchain] = stableTokens
	}
}

func GetStableCoins(blockchain Blockchain) []*Token {
	return stableCoins[blockchain]
}

func GetDefaultStableCoins(blockchain Blockchain) *Token {
	return GetStableCoins(blockchain)[0]
}

func GetToken(chain Blockchain, address string) *Token {
	for _, token := range tokens[chain] {
		if token.Address == address {
			return token
		}
	}
	return nil
}

func GetTokens(chain Blockchain) []*Token {
	return tokens[chain]
}
