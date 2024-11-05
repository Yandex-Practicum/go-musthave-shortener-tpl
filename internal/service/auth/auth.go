package auth

import (
	"fmt"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/errorscustom"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
)

// SecretSalt - соль для шифрования.
const (
	SecretSalt = "practicumSecretKey32"
	tokenSalt  = "tokenPracticum32"
)

// AuthService сервис отвечающий за авторизацию и верификацию.
//
//go:generate mockgen -source=auth.go -destination=mock_auth.go -package=auth
type AuthService interface {
	VerifyUser(token string) (string, error)
	CreatTokenForUser(userID string) (string, error)
}

// ServiceAuth - сервис для работы с JWT.
type ServiceAuth struct {
	passwordSalt []byte
	tokenSalt    []byte

	accessTokenTTL time.Duration
	userStorage    service.Storage
}

// NewServiceAuth - конструктор для сервиса авторизации.
func NewServiceAuth(storage service.Storage) *ServiceAuth {
	return &ServiceAuth{
		passwordSalt: []byte(SecretSalt),
		tokenSalt:    []byte(tokenSalt),

		accessTokenTTL: time.Hour,
		userStorage:    storage,
	}
}

// VerifyUser godoc
// @Tags AUTH_SERVICE
// @Summary Verify user
// @Description Verify user
// @Param VerifyUser body string true "token"
// @Success 200
// @Failure 500 "Internal server error"
// @Router / [options]
// VerifyUser проверяет наличие токена в заголовке Authorization.
func (sa *ServiceAuth) VerifyUser(token string) (string, error) {
	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorscustom.ErrBadVarifyToken
		}

		return sa.passwordSalt, nil
	})
	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("incorrect token: %v", err)
	}

	return claims.UserID, nil
}

// CreatTokenForUser создает JWT-токен для пользователя.
func (sa *ServiceAuth) CreatTokenForUser(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Устанавливаем срок действия токена на 24 часа
	})

	signedToken, err := token.SignedString(sa.passwordSalt)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
