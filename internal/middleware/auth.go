package middleware

import (
	"context"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
	"net/http"
)

type contextKey string

const (
	cookieName            = "token"
	secretKey  contextKey = "practicumSecretKey999"
)

type AuthMiddleware struct {
	authService *auth.ServiceAuth
}

func NewAuthMiddleware(authService *auth.ServiceAuth) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (a *AuthMiddleware) AuthMiddleware(h http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie(cookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userLogin, err := a.authService.VerifyUser(accessToken.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctxWithUser := context.WithValue(r.Context(), secretKey, userLogin)
		h.ServeHTTP(w, r.WithContext(ctxWithUser))
	}

	return http.HandlerFunc(fn)
}
