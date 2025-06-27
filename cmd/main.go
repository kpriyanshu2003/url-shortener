package main

import (
	"log"

	_ "github.com/kpriyanshu2003/url-shortener/docs"
	swagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kpriyanshu2003/url-shortener/internal/config"
	"github.com/kpriyanshu2003/url-shortener/internal/controller"
	"github.com/kpriyanshu2003/url-shortener/internal/database"
	"github.com/joho/godotenv"
)

func init() {
	if config.GetEnv("ENV", "development") == "development" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}


// @title           URL Shortener API
// @version         1.0
// @description     A simple URL shortener service
// @host            https://url-shortener-uy9w.onrender.com
// @BasePath        /

// @Summary      Render homepage
// @Description  Serves the index.html from public directory
// @Tags         Static
// @Produce      html
// @Success      200  {string}  string  "HTML page"
// @Router       / [get]
func main() {
	database.Init()
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())

	app.Static("/", "./public")

	app.Post("/", controller.ShortenURL)
	app.Get("/:code", controller.Redirect)
	app.Get("/swagger/*", swagger.WrapHandler)

	log.Fatal(app.Listen(":" + config.GetEnv("PORT", "3300")))
}
