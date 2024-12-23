package main

import (
	"encoding/json"
	"fmt"
	"github.com/holiman/uint256"
	"github.com/joho/godotenv"
	"gringotts/models"
	"gringotts/provider"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	if err := models.LoadBlockchain(); err != nil {
		log.Fatal(err)
	}
	if err := models.LoadTokens(); err != nil {
		log.Fatal(err)
	}
	if err := models.LoadAccounts(); err != nil {
		log.Fatal(err)
	}

	j := provider.Jupiter{}
	x, _ := j.GetSwap(&provider.SwapParams{
		Chain:       models.Solana,
		FromToken:   models.GetToken(models.Solana, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		ToToken:     models.GetToken(models.Solana, ""),
		Amount:      uint256.NewInt(2 * 1000 * 1000),
		Recipient:   "H6ieTjyWqcRFv1RwD9LghEFDxVtMaEMgVbgnYrdHMjr5",
		SlippageBPS: 100,
	})
	fmt.Println(x.Command)
	fmt.Println(x.Metadata)
	fmt.Println(x.OutAmount.String())
	fmt.Println(len(x.AddressLookup))
	z, _ := json.Marshal(x.Accounts)
	fmt.Println(string(z))

	app := fiber.New()
	registerRoutes(app)
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	godotenv.Load()
}
