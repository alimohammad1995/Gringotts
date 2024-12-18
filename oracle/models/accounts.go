package models

import (
	"encoding/json"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/spf13/viper"
	"gringotts/utils"
	"os"
)

const GringottsSeed = "Gringotts"
const PeerSeed = "Peer"

type Account struct {
	Address    string `json:"address"`
	IsWritable bool   `json:"is_writable"`
	IsSigner   bool   `json:"is_signer"`
}

type AccountModel struct {
	EstimateDiscriminator []byte     `json:"estimate_discriminator"`
	EstimateAccounts      []*Account `json:"estimate_accounts"`
	BridgeDiscriminator   []byte     `json:"bridge_discriminator"`
	PriceFeed             string     `json:"pyth_price_feed"`
}

var accountModel AccountModel

func LoadAccounts() error {
	content, err := os.ReadFile(viper.GetString("solana"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &accountModel); err != nil {
		return err
	}

	return nil
}

func GetEstimateAccounts() []*Account {
	return accountModel.EstimateAccounts
}

func GetEstimateDiscriminator() []byte {
	return accountModel.EstimateDiscriminator
}

func GetBridgeDiscriminator() []byte {
	return accountModel.BridgeDiscriminator
}

func GetPriceFeed() string {
	return accountModel.PriceFeed
}

func GetGringotts(chain Blockchain) string {
	return GetPDA(chain, [][]byte{[]byte(GringottsSeed)}).String()
}

func GetPeer(chain Blockchain, destination Blockchain) string {
	return GetPDA(chain, [][]byte{[]byte(PeerSeed), utils.ToLittleEndianBytes(destination.GetLzEId())}).String()
}

func GetSigner() string {
	return os.Getenv("SOLANA_SIGNER_SEED")
}

func GetPDA(chain Blockchain, seeds [][]byte) common.PublicKey {
	if chain != Solana && chain != SolanaDev {
		panic(fmt.Sprintf("need solana as the chain"))
	}

	pda, _, _ := common.FindProgramAddress(
		seeds,
		common.PublicKeyFromString(chain.GetContract()),
	)

	return pda
}
