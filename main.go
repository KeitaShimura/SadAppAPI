package main

import (
	"SadApp/src/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect();
	database.AutoMigrate();

	// Create a new Fiber app
	app := fiber.New()

	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Start the server
	log.Fatal(app.Listen(":8002"))
}
