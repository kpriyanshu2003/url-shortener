package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/database"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
	"github.com/kpriyanshu2003/url-shortener/internal/utils"
)

// ShortenURL godoc
// @Summary      Create a short URL
// @Description  Generates a short URL code for the given long URL. If no code is provided, one is generated.
// @Tags         URL
// @Accept       json
// @Produce      json
// @Param        payload  body      model.URLMapping  true  "Original URL and optional custom code"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       / [post]
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

// Redirect godoc
// @Summary      Redirect to original URL
// @Description  Redirects a short code to the corresponding original URL
// @Tags         URL
// @Produce      plain
// @Param        code  path      string  true  "Short code"
// @Success      302
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /{code} [get]
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
