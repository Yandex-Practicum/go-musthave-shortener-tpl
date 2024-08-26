package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"time"
)

const (
	SecretSalt = "practicumSecretKey32"
	tokenSalt  = "tokenPracticum32"
)

type AuthService interface {
	//CheckUserDB(userID string) error
	VerifyUser(token string) (string, error)
	CreatTokenForUser(userID string) (string, error)
}

type ServiceAuth struct {
	passwordSalt []byte
	tokenSalt    []byte

	accessTokenTTL time.Duration
	userStorage    service.Storage
}

func NewServiceAuth(storage service.Storage) *ServiceAuth {
	return &ServiceAuth{
		passwordSalt: []byte(SecretSalt),
		tokenSalt:    []byte(tokenSalt),

		accessTokenTTL: time.Hour,
		userStorage:    storage,
	}
}

//func (sa *ServiceAuth) CheckUserDB(userID string) error {
//	ok := sa.userStorage.CheckUser(userID)
//	if ok {
//		return nil
//	}
//	return errors.New("user does not exist")
//}

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

	return claims.UserID, nil
}

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
