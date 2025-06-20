package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/database"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
	"github.com/kpriyanshu2003/url-shortener/internal/utils"
)

func ShortenURL(c *fiber.Ctx) error {
	var payload model.URLMapping
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if payload.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}
	if payload.Code == "" {
		code, err := utils.GenerateUniqueCode()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate unique code",
			})
		}
		payload.Code = code
	}

	exists, err := database.CodeExists(payload.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Short URL code already exists",
		})
	}

	if err := database.InsertURL(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert URL into database",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"short_url":    "/" + payload.Code,
		"original_url": payload.URL,
	})
}

func Redirect(c *fiber.Ctx) error {
	code := c.Params("code")

	url, err := database.GetOriginalURL(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if url == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found",
		})
	}

	return c.Redirect(url, http.StatusTemporaryRedirect)
}
