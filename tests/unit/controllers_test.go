package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kpriyanshu2003/url-shortener/internal/controller"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
)

type MockDatabase struct {
	urls        map[string]string
	insertError bool
	existsError bool
	getError    bool
}

func (m *MockDatabase) CodeExists(code string) (bool, error) {
	if m.existsError {
		return false, errors.New("mock exists error")
	}
	_, exists := m.urls[code]
	return exists, nil
}

func (m *MockDatabase) InsertURL(mapping model.URLMapping) error {
	if m.insertError {
		return errors.New("mock insert error")
	}
	if _, exists := m.urls[mapping.Code]; exists {
		return errors.New("code already exists")
	}
	m.urls[mapping.Code] = mapping.URL
	return nil
}

func (m *MockDatabase) GetOriginalURL(code string) (string, error) {
	if m.getError {
		return "", errors.New("mock get error")
	}
	url, exists := m.urls[code]
	if !exists {
		return "", nil
	}
	return url, nil
}

type MockUtils struct {
	generateError bool
	fixedCode     string
}

func (m *MockUtils) GenerateUniqueCode() (string, error) {
	if m.generateError {
		return "", errors.New("mock generate error")
	}
	if m.fixedCode != "" {
		return m.fixedCode, nil
	}
	return "generated123", nil
}

func setupMocks() (*MockDatabase, *MockUtils) {
	return &MockDatabase{
		urls: make(map[string]string),
	}, &MockUtils{}
}

func TestShortenURLSuccess(t *testing.T) {

	mockDB, mockUtils := setupMocks()

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

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

func TestShortenURLWithGeneratedCode(t *testing.T) {

	mockDB, mockUtils := setupMocks()
	mockUtils.fixedCode = "auto456"

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

	urlMapping := model.URLMapping{
		URL: "https://example.com",
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

	if response["short_url"] != "/auto456" {
		t.Errorf("Expected short_url '/auto456', got '%s'", response["short_url"])
	}
}

func TestShortenURLInvalidJSON(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

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

func TestShortenURLEmptyURL(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

	urlMapping := model.URLMapping{
		URL:  "",
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "URL is required" {
		t.Errorf("Expected error 'URL is required', got '%s'", response["error"])
	}
}

func TestShortenURLCodeExists(t *testing.T) {
	mockDB, mockUtils := setupMocks()

	mockDB.urls["existing"] = "https://existing.com"

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

	urlMapping := model.URLMapping{
		URL:  "https://example.com",
		Code: "existing",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

func TestShortenURLGenerateCodeError(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	mockUtils.generateError = true

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

	urlMapping := model.URLMapping{
		URL: "https://example.com",
	}

	body, _ := json.Marshal(urlMapping)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestShortenURLDatabaseError(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	mockDB.existsError = true

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

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

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestShortenURLInsertError(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	mockDB.insertError = true

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Post("/", controller.ShortenURL)

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

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestRedirectSuccess(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	mockDB.urls["redirect123"] = "https://redirect-example.com"

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Get("/:code", controller.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/redirect123", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected status %d, got %d", http.StatusTemporaryRedirect, resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location != "https://redirect-example.com" {
		t.Errorf("Expected Location 'https://redirect-example.com', got '%s'", location)
	}
}

func TestRedirectNotFound(t *testing.T) {
	mockDB, mockUtils := setupMocks()

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Get("/:code", controller.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "Short URL not found" {
		t.Errorf("Expected error 'Short URL not found', got '%s'", response["error"])
	}
}

func TestRedirectDatabaseError(t *testing.T) {
	mockDB, mockUtils := setupMocks()
	mockDB.getError = true

	originalDB := controller.DB
	originalUtils := controller.Utils
	controller.DB = mockDB
	controller.Utils = mockUtils
	defer func() {
		controller.DB = originalDB
		controller.Utils = originalUtils
	}()

	app := fiber.New()
	app.Get("/:code", controller.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/test123", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "Database error" {
		t.Errorf("Expected error 'Database error', got '%s'", response["error"])
	}
}
