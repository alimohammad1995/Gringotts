package main

import (
	"encoding/binary"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gringotts/models"
	"log"
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

	programID := common.PublicKeyFromString("D4qgtRF5t7MzA5nASfag3ef5Y68xb5hBcp9ac9Wttfkr")

	seed1 := []byte("Peer")
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, 40231)

	// Derive the PDA
	pda, bump, err := common.FindProgramAddress(
		[][]byte{seed1, bytes},
		programID,
	)

	fmt.Println(pda, bump, err)

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
}
