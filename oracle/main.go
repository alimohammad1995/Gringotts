package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
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

	app := fiber.New()
	app.Use(cors.New())

	registerRoutes(app)

	if err := app.Listen("0.0.0.0:3000"); err != nil {
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
