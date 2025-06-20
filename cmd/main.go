package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/controller"
)

func main() {
	app := fiber.New()
	app.Static("/", "./public")

	app.Post("/", controller.ShortenURL)
	app.Get("/:code", controller.Redirect)

	log.Fatal(app.Listen(":3300"))
}
