package middleware

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	errIncorrect = errors.New("incorrect token")
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		login        string
		expectedErr  error
		token        string
		createErr    error
		expectedCode int
	}{
		{
			name:         "successful_auth",
			login:        "test",
			expectedCode: 200,
		},
		{
			name:         "bad_verify_user",
			expectedErr:  errIncorrect,
			createErr:    errIncorrect,
			expectedCode: 500,
		},
		{
			name:         "successful_creat_token",
			expectedErr:  errIncorrect,
			token:        "test",
			expectedCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Создаем запрос и ответ
			req := httptest.NewRequest("POST", "/", nil)
			resp := httptest.NewRecorder()

			// формируем handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("test_auth"))
			})

			// создаем заголовок
			resp.Header().Set("Authorization", tt.login)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := auth.NewMockAuthService(ctrl)
			mockAuth.EXPECT().VerifyUser(gomock.Any()).Return(tt.login, tt.expectedErr).AnyTimes()
			mockAuth.EXPECT().CreatTokenForUser(gomock.Any()).Return(tt.token, tt.createErr).AnyTimes()

			service := NewAuthMiddleware(mockAuth)

			// Применяем AuthMiddleware к хэндлеру
			wrappedHandler := service.AuthMiddleware(handler)

			// Вызываем middleware обернутый хэндлер
			wrappedHandler.ServeHTTP(resp, req)

			// Проверяем статус-код
			if status := resp.Code; status != tt.expectedCode {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
		})
	}

}

func TestCheckAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		expectedCookie bool
		expectedCode   int
		login          string
		expectedErr    error
	}{
		{
			name:         "successful_check",
			expectedCode: 200,
			login:        "test",
		},
		{
			name:         "incorrect_verify",
			expectedErr:  errIncorrect,
			expectedCode: 401,
		},
		{
			name:           "incorrect_cookie",
			expectedCookie: true,
			expectedCode:   401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest("GET", "/", nil)
			if !tt.expectedCookie {
				req.AddCookie(&http.Cookie{Name: string(UserIDContextKey), Value: "valid_token"})
			}

			resp := httptest.NewRecorder()

			mockAuth := auth.NewMockAuthService(ctrl)
			mockAuth.EXPECT().VerifyUser(gomock.Any()).Return(tt.login, tt.expectedErr).AnyTimes()

			service := NewAuthMiddleware(mockAuth)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})

			// Применяем AuthMiddleware к хэндлеру
			wrappedHandler := service.CheckAuthMiddleware(handler)

			// Вызываем middleware обернутый хэндлер
			wrappedHandler.ServeHTTP(resp, req)

			// Проверяем статус-код
			if status := resp.Code; status != tt.expectedCode {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
		})
	}
}
