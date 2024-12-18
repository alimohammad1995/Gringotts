package main

import (
	"github.com/gofiber/fiber/v2"
	"gringotts/handlers"
)

func registerRoutes(app *fiber.App) {
	app.Get("/transaction", handlers.HandleTransaction)
	app.Get("/network", handlers.HandleNetworks)
	app.Get("/network/:chain/token/", handlers.HandleTokens)
	app.Get("/test", handlers.HandleTest)
}
