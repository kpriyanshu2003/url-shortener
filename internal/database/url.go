package database

import (
	"database/sql"

	"github.com/kpriyanshu2003/url-shortener/internal/model"
)

const (
	_insertURL       = "INSERT INTO url_mappings (code, original) VALUES ($1, $2)"
	_getOriginalURL  = "SELECT original FROM url_mappings WHERE code = $1"
	_checkCodeExists = "SELECT EXISTS(SELECT 1 FROM url_mappings WHERE code = $1)"
)

func InsertURL(mapping model.URLMapping) error {
	_, err := DB.Exec(_insertURL, mapping.Code, mapping.URL)
	return err
}

func GetOriginalURL(code string) (string, error) {
	var original string
	err := DB.QueryRow(_getOriginalURL, code).Scan(&original)

	if err == sql.ErrNoRows {
		return "", nil
	}
	return original, err
}

func CodeExists(code string) (bool, error) {
	var exists bool
	err := DB.QueryRow(_checkCodeExists, code).Scan(&exists)
	return exists, err
}
