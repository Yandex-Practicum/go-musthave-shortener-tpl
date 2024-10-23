package models

import (
	"github.com/golang-jwt/jwt/v4"
)

// Claims - структура для хранения JWT-токена.
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

// User - структура для хранения данных пользователя.
type User struct {
	UUID  string `json:"UUID"`
	Token string `json:"Token"`
}

// UserURLs - структура для хранения URL пользователя.
type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
