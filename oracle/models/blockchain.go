package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
	"os"
	"slices"
)

type Blockchain string

type BlockchainModel struct {
	Chain    Blockchain
	Id       uint8  `json:"id"`
	LzEId    int    `json:"lz_eid"`
	Endpoint string `json:"endpoint"`
	Icon     string `json:"icon"`
	Contract string `json:"contract"`
}

const (
	Ethereum Blockchain = "ethereum"
	Binance  Blockchain = "binance"
	Polygon  Blockchain = "polygon"
	TRON     Blockchain = "tron"
	BASE     Blockchain = "base"
	Fantom   Blockchain = "fantom"
	Avax     Blockchain = "avax"
	Arbitrum Blockchain = "arbitrum"
	Optimism Blockchain = "optimism"
	Aurora   Blockchain = "aurora"
	Solana   Blockchain = "solana"

	_Solana   Blockchain = "_solana"
	_Arbitrum Blockchain = "_arbitrum"
)

var allValidBlockchains = []Blockchain{
	Ethereum,
	Binance,
	Polygon,
	TRON,
	BASE,
	Fantom,
	Avax,
	Arbitrum,
	Optimism,
	Aurora,
	Solana,

	_Solana,
	_Arbitrum,
}

var blockchains = map[Blockchain]*BlockchainModel{}

func LoadBlockchain() error {
	content, err := os.ReadFile(viper.GetString("blockchains"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &blockchains); err != nil {
		return err
	}

	idsMap := make(map[uint8]bool)
	lzEdiMap := make(map[int]bool)
	chainsMap := make(map[Blockchain]bool)

	for chain, info := range blockchains {
		if !slices.Contains(allValidBlockchains, chain) {
			return errors.New("blockchains not valid")
		}

		if idsMap[info.Id] || chainsMap[chain] || lzEdiMap[info.LzEId] {
			return errors.New(fmt.Sprintf("blockchain %s has duplicate data", chain))
		}

		idsMap[info.Id] = true
		chainsMap[chain] = true
		lzEdiMap[info.LzEId] = true
	}

	log.Info(fmt.Sprintf("Loaded %d blockchains", len(blockchains)))
	return nil
}

func GetAll() map[Blockchain]*BlockchainModel {
	return blockchains
}

func (b Blockchain) IsSupported() bool {
	return blockchains[b] != nil
}

func (b Blockchain) GetEndpoint() string {
	return blockchains[b].Endpoint
}

func (b Blockchain) GetId() uint8 {
	return blockchains[b].Id
}

func (b Blockchain) GetContract() string {
	return blockchains[b].Contract
}
