package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Open the connection to the database
	_, err := gorm.Open(mysql.Open("keita:@Keita8001@tcp(db:3306)/sadapp?parseTime=true"), &gorm.Config{})

	// Check for errors in opening the database connection
	if err != nil {
		panic("データベースに接続できませんでした。") // "Could not connect to the database."
	}

	// Create a new Fiber app
	app := fiber.New()

	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Start the server
	log.Fatal(app.Listen(":8002"))
}
