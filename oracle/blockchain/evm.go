package blockchain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"gringotts/models"
)

func GetConnection(blockchain models.Blockchain) (*ethclient.Client, error) {
	return ethclient.Dial(blockchain.GetEndpoint())
}
