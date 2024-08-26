package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
	"net/http"
)

const UserIDContextKey = "user_id"

type AuthMiddleware struct {
	authService *auth.ServiceAuth
}

func NewAuthMiddleware(authService *auth.ServiceAuth) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (a *AuthMiddleware) AuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		var accessToken string
		if authHeader != "" {
			accessToken = authHeader
		} else {
			//читаем токен из кук
			cookie, err := r.Cookie(UserIDContextKey)
			if err == nil && cookie.Value != "" {
				accessToken = cookie.Value
			} else {
				accessToken = ""
			}
		}

		userID, err := a.authService.VerifyUser(accessToken)
		if err != nil {

			// создаем токен
			userID = uuid.New().String()
			token, err := a.authService.CreatTokenForUser(userID)
			if err != nil {
				http.Error(w, `{"error":"Failed to generate auth token"}`, http.StatusInternalServerError)
				return
			}

			//создаем куку
			http.SetCookie(w, &http.Cookie{
				Name:     UserIDContextKey,
				Value:    token,
				HttpOnly: true,
			})

			// Устанавливаем заголовок Authorization
			w.Header().Set("Authorization", token)
		}
		ctxWithUser := context.WithValue(r.Context(), UserIDContextKey, userID)
		h.ServeHTTP(w, r.WithContext(ctxWithUser))
	}

	return http.HandlerFunc(fn)
}

func (a *AuthMiddleware) CheckAuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie(UserIDContextKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := a.authService.VerifyUser(accessToken.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), UserIDContextKey, userID)
		h.ServeHTTP(w, r.WithContext(ctxWithUser))
	}

	return http.HandlerFunc(fn)
}
