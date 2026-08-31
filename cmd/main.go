package main

import (
	"log"

	"github.com/ardhisparahita/cinema-booking-api/internal/config"
	"github.com/ardhisparahita/cinema-booking-api/pkg/database"
	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnv()

	app := fiber.New()

	database.ConnectDB()
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "body",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
