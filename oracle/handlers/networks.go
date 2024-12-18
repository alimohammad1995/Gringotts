package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gringotts/models"
)

type NetworkResponse struct {
	Chain models.Blockchain `json:"chain"`
	Name  string            `json:"name"`
	Icon  string            `json:"icon"`
}

func HandleNetworks(c *fiber.Ctx) error {
	res := make([]NetworkResponse, 0)

	for chain, blockchain := range models.GetAll() {
		res = append(res, NetworkResponse{
			Chain: chain,
			Name:  string(chain),
			Icon:  blockchain.Icon,
		})
	}

	return c.JSON(res)
}
