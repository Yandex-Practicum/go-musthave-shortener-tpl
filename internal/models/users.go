package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"l"`
}

type User struct {
	ID    string
	Token string
}
