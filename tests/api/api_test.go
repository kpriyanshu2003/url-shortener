package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/controller"
	"github.com/kpriyanshu2003/url-shortener/internal/database"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
)

func setupTestApp(t *testing.T) *fiber.App {
	database.Init()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)
	app.Get("/:code", controller.Redirect)

	return app
}

func TestAPIShortenURL(t *testing.T) {
	app := setupTestApp(t)

	urlMapping := model.URLMapping{
		URL:  "https://example.com",
		Code: "test123",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["short_url"] != "/test123" {
		t.Errorf("Expected short_url '/test123', got '%s'", response["short_url"])
	}
	if response["original_url"] != "https://example.com" {
		t.Errorf("Expected original_url 'https://example.com', got '%s'", response["original_url"])
	}
}

func TestAPIShortenURLWithoutCode(t *testing.T) {
	app := setupTestApp(t)

	urlMapping := model.URLMapping{
		URL: "https://google.com",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["original_url"] != "https://google.com" {
		t.Errorf("Expected original_url 'https://google.com', got '%s'", response["original_url"])
	}
	if response["short_url"] == "" {
		t.Error("Expected short_url to be generated, got empty string")
	}
}

func TestAPIShortenURLInvalidInput(t *testing.T) {
	app := setupTestApp(t)

	urlMapping := model.URLMapping{
		URL: "",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestAPIShortenURLDuplicateCode(t *testing.T) {
	app := setupTestApp(t)

	urlMapping := model.URLMapping{
		URL:  "https://first.com",
		Code: "duplicate",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform first request: %v", err)
	}
	resp.Body.Close()

	urlMapping2 := model.URLMapping{
		URL:  "https://second.com",
		Code: "duplicate",
	}

	body2, _ := json.Marshal(urlMapping2)
	req2 := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("Failed to perform second request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, resp2.StatusCode)
	}
}

func TestAPIRedirect(t *testing.T) {
	app := setupTestApp(t)

	urlMapping := model.URLMapping{
		URL:  "https://redirect-test.com",
		Code: "redirect123",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to create short URL: %v", err)
	}
	resp.Body.Close()

	req2 := httptest.NewRequest(http.MethodGet, "/redirect123", nil)
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("Failed to perform redirect request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected status %d, got %d", http.StatusTemporaryRedirect, resp2.StatusCode)
	}

	location := resp2.Header.Get("Location")
	if location != "https://redirect-test.com" {
		t.Errorf("Expected Location header 'https://redirect-test.com', got '%s'", location)
	}
}

func TestAPIRedirectNotFound(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestAPIInvalidJSON(t *testing.T) {
	app := setupTestApp(t)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
