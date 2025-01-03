package models

import (
	"encoding/json"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/spf13/viper"
	"gringotts/utils"
	"os"
)

const NativeMint = "So11111111111111111111111111111111111111112"
const GringottsSeed = "Gringotts"
const VaultSeed = "Vault"
const PeerSeed = "Peer"

type Account struct {
	Address    string `json:"address"`
	IsWritable bool   `json:"is_writable"`
	IsSigner   bool   `json:"is_signer"`
	Index      int    `json:"index"`
}

type AccountModel struct {
	EstimateDiscriminator []byte                    `json:"estimate_discriminator"`
	QuoteAccounts         map[Blockchain][]*Account `json:"quote_accounts"`
	SendAccounts          map[Blockchain][]*Account `json:"send_accounts"`
	BridgeDiscriminator   []byte                    `json:"bridge_discriminator"`
	PriceFeed             string                    `json:"pyth_price_feed"`
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
	for chain := range accountModel.QuoteAccounts {
		if !chain.IsSupported() {
			return fmt.Errorf("unsupported chain: %s", chain)
		}
	}
	for chain := range accountModel.SendAccounts {
		if !chain.IsSupported() {
			return fmt.Errorf("unsupported chain: %s", chain)
		}
	}

	return nil
}

func GetEstimateAccounts(chain Blockchain) []*Account {
	return accountModel.QuoteAccounts[chain]
}

func GetSendAccounts(chain Blockchain) []*Account {
	return accountModel.SendAccounts[chain]
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

func GetJupiter() string {
	return "JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"
}

func GetGringotts(chain Blockchain) string {
	return GetPDA(chain, [][]byte{[]byte(GringottsSeed)}).String()
}

func GetPeer(chain Blockchain, destination Blockchain) string {
	return GetPDA(chain, [][]byte{[]byte(PeerSeed), utils.ToBigEndianBytes(destination.GetLzEId())}).String()
}

func GetVault(chain Blockchain) string {
	return GetPDA(chain, [][]byte{[]byte(VaultSeed)}).String()
}

func GetAssociatedTokenAddress(owner string, token string) string {
	pda, _, _ := common.FindAssociatedTokenAddress(
		common.PublicKeyFromString(owner),
		common.PublicKeyFromString(token),
	)

	return pda.ToBase58()
}

func GetSigner() string {
	return os.Getenv("SOLANA_SIGNER_SEED")
}

func GetPDA(chain Blockchain, seeds [][]byte) common.PublicKey {
	pda, _, _ := common.FindProgramAddress(
		seeds,
		common.PublicKeyFromString(chain.GetContract()),
	)

	return pda
}
