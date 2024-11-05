package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGZipMiddleware(t *testing.T) {

	t.Run("testing_GZip_accept_encoding", func(t *testing.T) {
		// Создаем новый HTTP-запрос
		rRequest := httptest.NewRequest("POST", "/", nil)

		// Устанавливаем заголовок Accept-Encoding для GZip
		rRequest.Header.Set("Accept-Encoding", "gzip")

		// Создаем фиктивный HTTP-ответ
		wResponse := httptest.NewRecorder()

		// Создаем базовый хэндлер
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, world!"))
		})

		// Применяем GZip middleware к хэндлеру
		wrappedHandler := GZipMiddleware(handler)

		// Вызываем middleware обернутый хэндлер
		wrappedHandler.ServeHTTP(wResponse, rRequest)

		// Проверяем статус-код
		if status := wResponse.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
