package main

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	// описываем ожидаемое тело ответа при успешном запросе
	successBody := `{
     "response": {
         "text": "Извините, я пока ничего не умею"
     },
     "version": "1.0"
 }`

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		method       string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: successBody},
	}

	for _, tc := range testCases {
		t.Run("testing_"+tc.method, func(t *testing.T) {
			// создаём connect
			//dsm, _ := db.NewPstStorage("")

			r := httptest.NewRequest(tc.method, "/", nil)
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			logs := logger.NewLogger(logger.WithLevel("info"))
			// инициализируем файл для хранения
			fileName := "./test.txt"
			defer os.Remove(fileName)

			file, err := filestorage.NewSaveFile(fileName)
			if err != nil {
				logs.Error("Fatal", logger.ErrAttr(err))
			}
			defer file.Close()

			storage := mapstorage.NewMapURL()
			urlService := service.NewService(storage, logs)
			shortHandlers := handlers.NewHandlers(urlService, "http://localhost:8080/", logs, nil)

			switch tc.method {
			case http.MethodPost:
				shortHandlers.PostURL(w, r)
			case http.MethodGet:
				shortHandlers.GetURL(w, r)

			}

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				// assert.JSONEq помогает сравнить две JSON-строки
				assert.JSONEq(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
