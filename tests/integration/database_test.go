package database

import (
	"testing"

	"github.com/kpriyanshu2003/url-shortener/internal/database"
	"github.com/kpriyanshu2003/url-shortener/internal/model"
)

func setupTestDatabase(t *testing.T) {
	database.Init()
}

func TestInsertURL(t *testing.T) {
	setupTestDatabase(t)

	mapping := model.URLMapping{
		Code: "test001",
		URL:  "https://example.com/test",
	}

	err := database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert URL: %v", err)
	}

	exists, err := database.CodeExists(mapping.Code)
	if err != nil {
		t.Fatalf("Failed to check if code exists: %v", err)
	}
	if !exists {
		t.Error("Expected code to exist after insertion")
	}
}

func TestGetOriginalURL(t *testing.T) {
	setupTestDatabase(t)

	mapping := model.URLMapping{
		Code: "test002",
		URL:  "https://example.com/get-test",
	}

	err := database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert URL for test: %v", err)
	}

	retrievedURL, err := database.GetOriginalURL(mapping.Code)
	if err != nil {
		t.Fatalf("Failed to get original URL: %v", err)
	}
	if retrievedURL != mapping.URL {
		t.Errorf("Expected URL '%s', got '%s'", mapping.URL, retrievedURL)
	}
}

func TestGetOriginalURLNotFound(t *testing.T) {
	setupTestDatabase(t)

	retrievedURL, err := database.GetOriginalURL("nonexistent")
	if err != nil {
		t.Fatalf("Expected no error for non-existent code, got: %v", err)
	}
	if retrievedURL != "" {
		t.Errorf("Expected empty string for non-existent code, got '%s'", retrievedURL)
	}
}

func TestCodeExists(t *testing.T) {
	setupTestDatabase(t)

	exists, err := database.CodeExists("nonexistent123")
	if err != nil {
		t.Fatalf("Failed to check code existence: %v", err)
	}
	if exists {
		t.Error("Expected code to not exist")
	}

	mapping := model.URLMapping{
		Code: "test003",
		URL:  "https://example.com/exists-test",
	}

	err = database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert URL for test: %v", err)
	}

	exists, err = database.CodeExists(mapping.Code)
	if err != nil {
		t.Fatalf("Failed to check code existence: %v", err)
	}
	if !exists {
		t.Error("Expected code to exist")
	}
}

func TestInsertDuplicateURL(t *testing.T) {
	setupTestDatabase(t)

	mapping := model.URLMapping{
		Code: "duplicate001",
		URL:  "https://example.com/duplicate",
	}

	err := database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert URL first time: %v", err)
	}

	err = database.InsertURL(mapping)
	if err == nil {
		t.Error("Expected error when inserting duplicate code, got nil")
	}
}

func TestInsertAndRetrieveMultipleURLs(t *testing.T) {
	setupTestDatabase(t)

	mappings := []model.URLMapping{
		{Code: "multi001", URL: "https://example.com/1"},
		{Code: "multi002", URL: "https://example.com/2"},
		{Code: "multi003", URL: "https://example.com/3"},
	}

	for _, mapping := range mappings {
		err := database.InsertURL(mapping)
		if err != nil {
			t.Fatalf("Failed to insert URL %s: %v", mapping.Code, err)
		}
	}

	for _, mapping := range mappings {
		retrievedURL, err := database.GetOriginalURL(mapping.Code)
		if err != nil {
			t.Fatalf("Failed to get original URL for %s: %v", mapping.Code, err)
		}
		if retrievedURL != mapping.URL {
			t.Errorf("Expected URL '%s' for code '%s', got '%s'",
				mapping.URL, mapping.Code, retrievedURL)
		}

		exists, err := database.CodeExists(mapping.Code)
		if err != nil {
			t.Fatalf("Failed to check existence of code %s: %v", mapping.Code, err)
		}
		if !exists {
			t.Errorf("Expected code '%s' to exist", mapping.Code)
		}
	}
}

func TestEmptyCodeHandling(t *testing.T) {
	setupTestDatabase(t)

	mapping := model.URLMapping{
		Code: "",
		URL:  "https://example.com/empty-code",
	}

	err := database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert URL with empty code: %v", err)
	}

	retrievedURL, err := database.GetOriginalURL("")
	if err != nil {
		t.Fatalf("Failed to get URL with empty code: %v", err)
	}
	if retrievedURL != mapping.URL {
		t.Errorf("Expected URL '%s' for empty code, got '%s'", mapping.URL, retrievedURL)
	}
}

func TestSpecialCharactersInURL(t *testing.T) {
	setupTestDatabase(t)

	mapping := model.URLMapping{
		Code: "special001",
		URL:  "https://example.com/path?query=value&other=123#fragment",
	}

	err := database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert URL with special characters: %v", err)
	}

	retrievedURL, err := database.GetOriginalURL(mapping.Code)
	if err != nil {
		t.Fatalf("Failed to get URL with special characters: %v", err)
	}
	if retrievedURL != mapping.URL {
		t.Errorf("Expected URL '%s', got '%s'", mapping.URL, retrievedURL)
	}
}

func TestLongURL(t *testing.T) {
	setupTestDatabase(t)

	longURL := "https://example.com/very-long-path/" +
		"this-is-a-very-long-url-segment-that-might-test-database-limits/" +
		"another-long-segment-with-many-characters-to-ensure-proper-handling/" +
		"final-segment?query=very-long-query-parameter-value&another=parameter"

	mapping := model.URLMapping{
		Code: "long001",
		URL:  longURL,
	}

	err := database.InsertURL(mapping)
	if err != nil {
		t.Fatalf("Failed to insert long URL: %v", err)
	}

	retrievedURL, err := database.GetOriginalURL(mapping.Code)
	if err != nil {
		t.Fatalf("Failed to get long URL: %v", err)
	}
	if retrievedURL != mapping.URL {
		t.Errorf("Expected long URL to match, got different URL")
	}
}
