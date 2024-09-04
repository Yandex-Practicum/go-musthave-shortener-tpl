package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

type User struct {
	UUID  string `json:"UUID"`
	Token string `json:"Token"`
}

type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
