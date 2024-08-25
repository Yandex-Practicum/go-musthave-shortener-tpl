package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/db"
	"time"
)

type ServiceAuth struct {
	passwordSalt []byte
	tokenSalt    []byte

	accessTokenTTL time.Duration
	userStorage    *db.PstStorage
}

func NewServiceAuth(passwordSalt, tokenSalt []byte, storage *db.PstStorage) *ServiceAuth {
	return &ServiceAuth{
		passwordSalt: passwordSalt,
		tokenSalt:    tokenSalt,

		accessTokenTTL: time.Hour,
		userStorage:    storage,
	}
}

func (sa *ServiceAuth) VerifyUser(token string) (string, error) {
	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return sa.tokenSalt, nil
	})
	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("incorrect token: %v", err)
	}

	return claims.Login, nil
}
