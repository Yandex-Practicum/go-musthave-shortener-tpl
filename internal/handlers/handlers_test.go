package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/mocks"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Предполагаем, что функция EncodeURL и переменная MapStorage уже определены в вашем пакете

func TestPostURL(t *testing.T) {
	// Тест на успешное кодирование URL
	logs := logger.NewLogger(logger.WithLevel("info"))
	storage := mapstorage.NewMapURL()

	urlService := service.NewService(storage, logs)
	shortHandlers := NewHandlers(urlService, "http://localhost:8080", logs, nil)

	t.Run("test_post_URL", func(t *testing.T) {
		payload := []byte("http://example.com")
		rRequest := httptest.NewRequest("POST", "/url", bytes.NewBuffer(payload))
		wResonse := httptest.NewRecorder()

		shortHandlers.PostURL(wResonse, rRequest)

		// Проверяем, что статус ответа - 201 Created
		assert.Equal(t, http.StatusCreated, wResonse.Code)

		// Проверяем, что тело ответа содержит URL
		responseURL := wResonse.Body.String()
		assert.Contains(t, responseURL, "http://localhost:8080/")

		// Проверяем, что в MapStorage добавлен новый URL
		encodedURL := strings.TrimPrefix(responseURL, "http://localhost:8080/")
		originalURL, err := storage.GetURL(encodedURL)
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com", originalURL)
	})

	//Тест на обработку пустого тела запроса
	t.Run("empty_request_body", func(t *testing.T) {
		rRequest := httptest.NewRequest("POST", "/url", bytes.NewBuffer([]byte("")))
		wResonse := httptest.NewRecorder()

		shortHandlers.PostURL(wResonse, rRequest)

		// Проверяем, что статус ответа - 200 OK
		assert.Equal(t, http.StatusNotFound, wResonse.Code)
	})
}

func TestHandlersPostJSON(t *testing.T) {
	logs := logger.NewLogger(logger.WithLevel("info"))
	storage := mapstorage.NewMapURL()

	urlService := service.NewService(storage, logs)
	shortHandlers := NewHandlers(urlService, "http://localhost:8080", logs, nil)

	t.Run("test_post_JSON", func(t *testing.T) {
		payload := "{\"url\": \"https://practicum.yandex.ru\"}"
		param := strings.NewReader(payload)
		rRequest := httptest.NewRequest("POST", "/", param)
		wResonse := httptest.NewRecorder()

		shortHandlers.PostJSON(wResonse, rRequest)

		// Проверяем, что статус ответа - 201 Created
		assert.Equal(t, http.StatusCreated, wResonse.Code)

	})

	t.Run("test_post_JSON_noBody", func(t *testing.T) {
		payload := ""
		param := strings.NewReader(payload)
		rRequest := httptest.NewRequest("POST", "/", param)
		wResonse := httptest.NewRecorder()

		shortHandlers.PostJSON(wResonse, rRequest)

		//проверяем пустое тело
		assert.Equal(t, http.StatusNotFound, wResonse.Code)
	})
}

func TestGetURL(t *testing.T) {
	// Тест на успешное декодирование URL
	logs := logger.NewLogger(logger.WithLevel("info"))

	storage := mapstorage.NewMapURL()

	urlService := service.NewService(storage, logs)
	shortHandlers := NewHandlers(urlService, "http://localhost:8080", logs, nil)
	t.Run("test_get_URL", func(t *testing.T) {

		payload := []byte("http://example.com")
		rRequest := httptest.NewRequest("POST", "/url", bytes.NewBuffer(payload))
		wResonse := httptest.NewRecorder()

		shortHandlers.PostURL(wResonse, rRequest)

		responseURL := wResonse.Body.String()
		encodedURL := strings.TrimPrefix(responseURL, "http://localhost:8080/")
		rRequest = httptest.NewRequest("GET", "http://localhost:8080/", nil)
		wResonse = httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		rRequest = rRequest.WithContext(context.WithValue(rRequest.Context(), chi.RouteCtxKey, chiCtx))
		chiCtx.URLParams.Add("id", encodedURL)

		shortHandlers.GetURL(wResonse, rRequest)

		// Проверяем, что статус ответа - 200 OK
		assert.Equal(t, http.StatusTemporaryRedirect, wResonse.Code)

		// Проверяем, что в MapStorage добавлен новый URL
		originalURL, err := storage.GetURL(encodedURL)
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com", originalURL)

		// Проверяем, что в MapStorage нет URL
		rRequest = httptest.NewRequest("GET", "http://localhost:8080/", nil)
		chiCtx = chi.NewRouteContext()
		rRequest = rRequest.WithContext(context.WithValue(rRequest.Context(), chi.RouteCtxKey, chiCtx))
		chiCtx.URLParams.Add("id", "Nourl")
		wResonse = httptest.NewRecorder()
		shortHandlers.GetURL(wResonse, rRequest)
		assert.Equal(t, http.StatusNotFound, wResonse.Code)
	})
}

func TestGetPing(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgre := mocks.NewMockStorage(ctrl)
	mockPostgre.EXPECT().Ping().Return(nil)

	logger := logger.NewLogger(logger.WithLevel("info"))

	service := service.NewService(mockPostgre, logger)
	handlers := &Handlers{service: service, baseURL: "http://localhost:8080/", logger: logger}

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handlers.GetPing(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ожидался статус %d, но получен %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("ожидался заголовок %s, но получен %s", "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	}
}

//func TestPostBatchDB_Success(t *testing.T) {
//	testURL := "https://youtube.com"
//
//	multipleURL := models.MultipleURL{
//		CorrelationID: "1",
//		OriginalURL:   testURL,
//	}
//
//	marsh, err := json.Marshal([]models.MultipleURL{multipleURL})
//	if err != nil {
//		t.Fatal(err)
//	}
//	buf := bytes.NewReader(marsh)
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockPostgre := mocks.NewMockStorage(ctrl)
//	mockPostgre.EXPECT().CheckURL(testURL).Return("", nil)
//	mockPostgre.EXPECT().SaveURL(gomock.Any(), testURL).Return(nil)
//
//	logger := logger.NewLogger(logger.WithLevel("info"))
//
//	service := service.NewService(mockPostgre, logger)
//	handlers := &Handlers{service: service, baseURL: "http://localhost:8080/", logger: logger}
//
//	req, err := http.NewRequest("POST", "/api/shorten/batch", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	w := httptest.NewRecorder()
//	handlers.PostBatchDB(w, req)
//
//	if w.Code != http.StatusCreated {
//		t.Errorf("ожидался статус %d, но получен %d", http.StatusCreated, w.Code)
//	}
//}

func TestPostBatchDB_InCorrectRequest(t *testing.T) {

	buf := bytes.NewReader([]byte("lololo"))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgre := mocks.NewMockStorage(ctrl)

	logger := logger.NewLogger(logger.WithLevel("info"))

	service := service.NewService(mockPostgre, logger)
	handlers := &Handlers{service: service, baseURL: "http://localhost:8080/", logger: logger}

	req, err := http.NewRequest("POST", "/api/shorten/batch", buf)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handlers.PostBatchDB(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус %d, но получен %d", http.StatusBadRequest, w.Code)
	}
}

//func TestPostBatchDB_StorageError(t *testing.T) {
//	testURL := "https://youtube.com"
//
//	multipleURL := models.MultipleURL{
//		CorrelationID: "1",
//		OriginalURL:   testURL,
//	}
//
//	marsh, err := json.Marshal([]models.MultipleURL{multipleURL})
//	if err != nil {
//		t.Fatal(err)
//	}
//	buf := bytes.NewReader(marsh)
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockPostgre := mocks.NewMockStorage(ctrl)
//	mockErr := errors.New("Some error")
//	mockPostgre.EXPECT().CheckURL(testURL).Return("", nil)
//	mockPostgre.EXPECT().SaveURL(gomock.Any(), testURL).Return(mockErr)
//
//	logger := logger.NewLogger(logger.WithLevel("info"))
//
//	service := service.NewService(mockPostgre, logger)
//	handlers := &Handlers{service: service, baseURL: "http://localhost:8080/", logger: logger}
//
//	req, err := http.NewRequest("POST", "/api/shorten/batch", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	w := httptest.NewRecorder()
//	handlers.PostBatchDB(w, req)
//
//	if w.Code != http.StatusInternalServerError {
//		t.Errorf("ожидался статус %d, но получен %d", http.StatusInternalServerError, w.Code)
//	}
//}

func TestPostBatchDB_EmptyRequest(t *testing.T) {

	buf := bytes.NewReader([]byte(""))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgre := mocks.NewMockStorage(ctrl)

	logger := logger.NewLogger(logger.WithLevel("info"))

	service := service.NewService(mockPostgre, logger)
	handlers := &Handlers{service: service, baseURL: "http://localhost:8080/", logger: logger}

	req, err := http.NewRequest("POST", "/api/shorten/batch", buf)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handlers.PostBatchDB(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("ожидался статус %d, но получен %d", http.StatusNotFound, w.Code)
	}
}
