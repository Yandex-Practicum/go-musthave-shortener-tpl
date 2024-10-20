package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/workers"
	"net/http"
	//_ "net/http/pprof"
)

func main() {
	// инициализируем конфиг
	configs := NewConfigs()
	configs.Parse()

	// инициализируем logger
	logs := logger.NewLogger(logger.WithLevel(configs.LogLevel))
	logs.Info("Start logger")

	// инициализируем хранилище
	repo := initDB(configs.AddrDB, configs.PathFile)
	logs.Info("Connecting DB")
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

	// инициализируем worker
	worker := workers.NewWorkerDeleted(urlService)

	// передаем в хенлер сервис и baseURL
	shortHandlers := handlers.NewHandlers(urlService, configs.BaseURL, logs, worker)
	logs.Info(fmt.Sprintf("Handlers created PORT: %s", configs.AddrServer))

	// инициализировали роутер и создали Post и Get
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)

	// Добавляем pprof маршруты вручную
	r.Mount("/debug", middleware2.Profiler())

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.GZipMiddleware)
		r.Use(authorization.AuthMiddleware)
		r.Post("/", shortHandlers.PostURL)
		r.Post("/api/shorten", shortHandlers.PostJSON)
		r.Post("/api/shorten/batch", shortHandlers.PostBatchDB)
	})

	r.Get("/{id}", shortHandlers.GetURL)
	r.Get("/ping", shortHandlers.GetPing)

	r.Route("/api/user/urls", func(r chi.Router) {
		r.Use(authorization.CheckAuthMiddleware)
		r.Get("/", shortHandlers.GetUsersURLs)
		r.Delete("/", shortHandlers.DeletionURLs)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.StartWorkerDeletion(ctx)

	// слушаем выбранны порт = configs.AddrServer
	if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
		logs.Error("Err:", logger.ErrAttr(err))
	}
}
