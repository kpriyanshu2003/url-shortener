package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/database"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
	"github.com/kpriyanshu2003/url-shortener/internal/utils"
)

type DatabaseInterface interface {
	CodeExists(code string) (bool, error)
	InsertURL(mapping model.URLMapping) error
	GetOriginalURL(code string) (string, error)
}

type UtilsInterface interface {
	GenerateUniqueCode() (string, error)
}

type RealDatabase struct{}

func (r *RealDatabase) CodeExists(code string) (bool, error) {
	return database.CodeExists(code)
}

func (r *RealDatabase) InsertURL(mapping model.URLMapping) error {
	return database.InsertURL(mapping)
}

func (r *RealDatabase) GetOriginalURL(code string) (string, error) {
	return database.GetOriginalURL(code)
}

type RealUtils struct{}

func (r *RealUtils) GenerateUniqueCode() (string, error) {
	return utils.GenerateUniqueCode()
}

var (
	DB    DatabaseInterface = &RealDatabase{}
	Utils UtilsInterface    = &RealUtils{}
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
		code, err := Utils.GenerateUniqueCode()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate unique code",
			})
		}
		payload.Code = code
	}

	exists, err := DB.CodeExists(payload.Code)
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

	if err := DB.InsertURL(payload); err != nil {
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

	url, err := DB.GetOriginalURL(code)
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
