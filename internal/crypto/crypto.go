package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

const (
	cookieName = "user_id"
	secretKey  = "practicumSecretKey999"
)

// Генерация подписи для куки
func GenSign(data string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// Проверка подписи куки
func CheckSign(data, signature string) bool {
	expectedSignature := GenSign(data)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// Генерация нового уникального идентификатора пользователя
func GenNewUserID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}

// Установка куки с уникальным ID пользователя
func setUserCookie(w http.ResponseWriter, userID string) {

	// Подписываем ID пользователя
	signedUserID := fmt.Sprintf("%s|%s", userID, GenSign(userID))

	// Создаём и устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    signedUserID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}
