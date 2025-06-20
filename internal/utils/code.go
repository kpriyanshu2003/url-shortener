package utils

import (
	"math/rand"
	"time"

	"github.com/kpriyanshu2003/url-shortener/internal/database"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomCode() string {
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateUniqueCode() (string, error) {
	for {
		code := generateRandomCode()
		exists, err := database.CodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
}
