package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Routes
	// app.Post("/shorten", handler.ShortenURL)
	// app.Get("/:code", handler.Redirect)

	log.Fatal(app.Listen(":3000"))
}
