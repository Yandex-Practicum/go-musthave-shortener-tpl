package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"net/http"
)

func main() {
	// инициализируем конфиг
	configs := NewConfigs()
	configs.Parse()

	// инициализируем logger
	logs := logger.NewLogger(logger.WithLevel(configs.LogLevel))
	logs.Info("Start logger")

	// инициализируем хранилище
	repo := configs.repository
	logs.Info("Connecting DB", repo)
	defer repo.Close()

	// инициализируем сервис
	urlService := service.NewService(repo, logs)
	logs.Info(("Service created"))

	// передаем в хенлер сервис и baseURL
	shortHandlers := handlers.NewHandlers(urlService, configs.BaseURL, logs)
	logs.Info(fmt.Sprintf("Handlers created PORT: %s", configs.AddrServer))

	// инициализировали роутер и создали Post и Get
	r := chi.NewRouter()
	r.Post("/", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostURL)))
	r.Post("/api/shorten", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostJSON)))
	r.Get("/{id}", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.GetURL)))
	r.Get("/ping", middleware.WithLogging(shortHandlers.GetPing))

	// слушаем выбранны порт = configs.AddrServer
	if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
		logs.Error("Err:", logger.ErrAttr(err))
	}
}
