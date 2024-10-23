package handlers

import (
	"bytes"
	"context"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	serv "github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"net/http/httptest"
)

func ExampleHandlers_PostURL() {

	// Создаем  хранилище
	storage := mapstorage.NewMapURL()

	// Создаем логгер
	loger := logger.NewLogger(logger.WithLevel("info"))
	service := serv.NewService(storage, loger)
	baseURL := "http://localhost:8080"

	h := &Handlers{
		service: service,
		baseURL: baseURL,
		logger:  loger,
	}
	payload := []byte("http://example.com")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))
	ctx := context.WithValue(r.Context(), middleware.UserIDContextKey, "userID")
	r = r.WithContext(ctx)

	// Test
	h.PostURL(w, r)

	// Output:
	// "TmMIw"
}
