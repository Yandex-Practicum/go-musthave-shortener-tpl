package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
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
	storage := mapstorage.NewMapURL()
	urlService := service.NewService(storage)
	shortHandlers := NewHandlers(urlService, "http://localhost:8080")

	t.Run("Success", func(t *testing.T) {
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
	t.Run("EmptyRequestBody", func(t *testing.T) {
		rRequest := httptest.NewRequest("POST", "/url", bytes.NewBuffer([]byte("")))
		wResonse := httptest.NewRecorder()

		shortHandlers.PostURL(wResonse, rRequest)

		// Проверяем, что статус ответа - 200 OK
		assert.Equal(t, http.StatusOK, wResonse.Code)
	})
}

func TestGetURL(t *testing.T) {
	// Тест на успешное декодирование URL
	storage := mapstorage.NewMapURL()
	urlService := service.NewService(storage)
	shortHandlers := NewHandlers(urlService, "http://localhost:8080")
	t.Run("Success", func(t *testing.T) {

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
