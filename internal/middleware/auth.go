package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

type AuthMiddleware struct {
	authService auth.AuthService
}

func NewAuthMiddleware(authService auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func GetUserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(UserIDContextKey).(models.User)
	return user, ok
}

// AuthMiddleware godoc
// @Tags MIDDLEWARE
// @Summary Auth middleware
// @Description Auth middleware
// @Param Authorization header string true "Authorization"
// @Success 200
// @Failure 500 "Internal server error"
// @Router / [PATCH]
// AuthMiddleware проверяет наличие токена в заголовке Authorization и создает куку с токеном
func (a *AuthMiddleware) AuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		var accessToken string
		if authHeader != "" {
			accessToken = authHeader
		} else {
			//читаем токен из кук
			cookie, err := r.Cookie(string(UserIDContextKey))
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
				Name:     string(UserIDContextKey),
				Value:    token,
				HttpOnly: true,
				Path:     "/",
			})

			// Устанавливаем заголовок Authorization
			w.Header().Set("Authorization", token)
		}
		ctxWithUser := context.WithValue(r.Context(), UserIDContextKey, userID)
		h.ServeHTTP(w, r.WithContext(ctxWithUser))
	}

	return http.HandlerFunc(fn)
}

// CheckAuthMiddleware godoc
// @Tags MIDDLEWARE
// @Summary Check auth middleware
// @Description Check auth middleware
// @Param CheckAuthMiddleware header string true "user_id"
// @Success 200
// @Failure 401 "Unauthorized"
// @Router /api/user/urls [PATCH]
// CheckAuthMiddleware проверяет наличие токена в куке
func (a *AuthMiddleware) CheckAuthMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie(string(UserIDContextKey))
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
