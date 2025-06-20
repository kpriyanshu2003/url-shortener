package database

import (
	"database/sql"

	"github.com/kpriyanshu2003/url-shortener/internal/model"
)

const (
	_insertURL       = "INSERT INTO url_shortener (code, url) VALUES ($1, $2)"
	_getOriginalURL  = "SELECT url FROM url_shortener WHERE code = $1"
	_checkCodeExists = "SELECT EXISTS(SELECT 1 FROM url_shortener WHERE code = $1)"
)

func InsertURL(mapping model.URLMapping) error {
	_, err := DB.Exec(_insertURL, mapping.Code, mapping.URL)
	return err
}

func GetOriginalURL(code string) (string, error) {
	var url string
	err := DB.QueryRow(_getOriginalURL, code).Scan(&url)

	if err == sql.ErrNoRows {
		return "", nil
	}
	return url, err
}

func CodeExists(code string) (bool, error) {
	var exists bool
	err := DB.QueryRow(_checkCodeExists, code).Scan(&exists)
	return exists, err
}
