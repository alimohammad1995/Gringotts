package blockchain

import (
	"context"
	"github.com/blocto/solana-go-sdk/client"
	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/address_lookup_table"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/gofiber/fiber/v2/log"
	"gringotts/models"
)

var altCache map[string]types.AddressLookupTableAccount

func GetALT(chain models.Blockchain, altsAddress []string) []types.AddressLookupTableAccount {
	alts := make([]types.AddressLookupTableAccount, 0, len(altsAddress))

	for _, altAddress := range altsAddress {
		alt, err := getALT(chain, altAddress)
		if err != nil {
			log.Error(err)
			continue
		}
		
		alts = append(alts, alt)
	}

	return alts
}

func getALT(chain models.Blockchain, altAddress string) (types.AddressLookupTableAccount, error) {
	if alt, found := altCache[altAddress]; found {
		return alt, nil
	}

	accountInfo, err := GetSOLConnection(chain).GetAccountInfo(context.Background(), altAddress)
	if err != nil {
		return types.AddressLookupTableAccount{}, err
	}

	addressLookupTable, err := address_lookup_table.DeserializeLookupTable(accountInfo.Data, accountInfo.Owner)
	if err != nil {
		log.Errorw("invalid altAddress", "err", err)
		return types.AddressLookupTableAccount{}, err
	}

	altCache[altAddress] = types.AddressLookupTableAccount{
		Key:       solana.PublicKeyFromString(altAddress),
		Addresses: addressLookupTable.Addresses,
	}

	return altCache[altAddress], nil
}

func GetSOLConnection(chain models.Blockchain) *client.Client {
	return client.NewClient(chain.GetEndpoint())
}
