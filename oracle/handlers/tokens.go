package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gringotts/models"
)

type TokenResponse struct {
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	Icon     string `json:"icon"`
	Address  string `json:"address"`
}

func HandleTokens(c *fiber.Ctx) error {
	chain := models.Blockchain(c.Params("chain"))

	if !chain.IsSupported() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("%s is not a supported blockchain", chain)})
	}

	res := make([]TokenResponse, 0)

	for _, token := range models.GetTokens(chain) {
		res = append(res, TokenResponse{
			Name:     token.Name,
			Decimals: token.Decimals,
			Icon:     token.Icon,
			Address:  token.Address,
		})
	}

	return c.JSON(res)
}
