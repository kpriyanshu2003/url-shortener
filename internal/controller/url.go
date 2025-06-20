package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
)

var urlStore = make(map[string]string) // temp store

func ShortenURL(c *fiber.Ctx) error {
	var payload model.URLMapping
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if payload.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Iriginal URL is required",
		})
	}
	if payload.Code == "" {
		// Generate a random code if not provided
		payload.Code = "short" // This should be replaced with a proper random code generation logic
	}

	// DB interaction
	urlStore[payload.Code] = payload.URL

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"short_url":    "/" + payload.Code,
		"original_url": payload.URL,
	})
}

func Redirect(c *fiber.Ctx) error {
	code := c.Params("code")

	// DB interaction
	originalURL, found := urlStore[code]
	if !found {
		return c.Status(fiber.StatusNotFound).SendString("Short URL not found")
	}
	return c.Redirect(originalURL, http.StatusTemporaryRedirect)
}
