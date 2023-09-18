package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
