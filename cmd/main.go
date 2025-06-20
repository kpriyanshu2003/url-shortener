package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kpriyanshu2003/url-shortener/internal/controller"
	"github.com/kpriyanshu2003/url-shortener/internal/utils"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())

	app.Static("/", "./public")

	app.Post("/", controller.ShortenURL)
	app.Get("/:code", controller.Redirect)

	log.Fatal(app.Listen(":" + utils.GetEnv("PORT", "3300")))
}
