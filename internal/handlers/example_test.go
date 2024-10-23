package handlers

import (
	"bytes"
	"context"
	"fmt"
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

	// создаем body
	payload := []byte("http://example.com")

	// создаем респонс
	w := httptest.NewRecorder()

	// создаем запрос
	r := httptest.NewRequest("POST", baseURL, bytes.NewBuffer(payload))

	// создаем контекст
	ctx := context.WithValue(r.Context(), middleware.UserIDContextKey, "userID")
	r = r.WithContext(ctx)

	// Test
	h.PostURL(w, r)

	// Print
	fmt.Println(w.Code)

	// Output:
	// 201
}
