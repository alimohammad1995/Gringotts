package provider

import (
	"gringotts/models"
)

type Provider interface {
	GetSwap(params *models.SwapParams) (*models.Swap, error)
}
