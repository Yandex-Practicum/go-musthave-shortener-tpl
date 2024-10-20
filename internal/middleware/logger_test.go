package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithLogging(t *testing.T) {
	t.Run("testing_middleware_logger", func(t *testing.T) {

		// Создаем новый HTTP-запрос
		req := httptest.NewRequest("POST", "/", nil)

		// Создаем фиктивный HTTP-ответ
		resp := httptest.NewRecorder()

		// Создаем базовый хэндлер
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test"))
		})

		// Применяем Logger middleware к хэндлеру
		wrappedHandler := WithLogging(handler)

		// Вызываем middleware обернутый хэндлер
		wrappedHandler.ServeHTTP(resp, req)

		// Проверяем статус-код
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

}
