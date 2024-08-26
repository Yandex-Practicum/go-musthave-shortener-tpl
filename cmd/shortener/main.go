package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
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

	urlService := service.NewService(
		repo,
		logs,
	)
	logs.Info(("Service created"))

	// инициализируем проверку авторизацию
	serviceAuth := auth.NewServiceAuth(repo)
	authorization := middleware.NewAuthMiddleware(serviceAuth)

	// передаем в хенлер сервис и baseURL
	shortHandlers := handlers.NewHandlers(urlService, configs.BaseURL, logs)
	logs.Info(fmt.Sprintf("Handlers created PORT: %s", configs.AddrServer))

	// инициализировали роутер и создали Post и Get
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Group(func(r chi.Router) {
		r.Use(middleware.GZipMiddleware)
		r.Use(authorization.AuthMiddleware)
		r.Post("/", shortHandlers.PostURL)
		r.Post("/api/shorten", shortHandlers.PostJSON)
		r.Post("/api/shorten/batch", shortHandlers.PostBatchDB)
	})

	r.Get("/{id}", shortHandlers.GetURL)
	r.Get("/ping", shortHandlers.GetPing)

	r.Group(func(r chi.Router) {
		r.Use(authorization.CheckAuthMiddleware)
		r.Get("/api/user/urls", shortHandlers.GetUsersURLs)
	})

	//r.Post("/", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostURL)))
	//r.Post("/api/shorten", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostJSON)))
	//r.Post("/api/shorten/batch", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.PostBatchDB)))
	//r.Get("/{id}", middleware.WithLogging(middleware.GZipMiddleware(shortHandlers.GetURL)))
	//r.Get("/ping", middleware.WithLogging(shortHandlers.GetPing))

	// слушаем выбранны порт = configs.AddrServer
	if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
		logs.Error("Err:", logger.ErrAttr(err))
	}
}
