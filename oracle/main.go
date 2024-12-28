package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gringotts/models"
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
