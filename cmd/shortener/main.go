package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/fileStorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"net/http"
)

func main() {
	// инициализируем конфиг
	configs := NewConfigs()
	configs.Parse()

	// инициализируем logger
	logs := logger.NewLogger(logger.WithLevel(configs.flagLogLevel))
	logs.Info("Start logger")

	// инициализируем хранилище
	storage := mapstorage.NewMapURL()
	logs.Info("Storage created")

	// инициализируем файл для хранения
	file, err := fileStorage.NewSaveFile(configs.flagPathDB)
	if err != nil {
		logs.Error("Fatal", logger.ErrAttr(err))
	}
	defer file.Close()

	// передаем в сервис хранилище
	urlService := service.NewService(storage, file)
	logs.Info(("Service created"))

	// передаем в хенлер сервис и baseURL
	shortHandlers := handlers.NewHandlers(urlService, configs.BaseURL, logs)
	logs.Info(fmt.Sprintf("Handlers created PORT: %s", configs.AddrServer))

	// инициализировали роутер и создали Post и Get
	r := chi.NewRouter()
	r.Post("/", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostURL)))
	r.Post("/api/shorten", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostJSON)))
	r.Get("/{id}", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.GetURL)))

	// слушаем выбранны порт = configs.AddrServer
	if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
		logs.Error("Err:", logger.ErrAttr(err))
	}
}
