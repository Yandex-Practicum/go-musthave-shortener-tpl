package main

import (
	"context"
	"fmt"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"net/http"

	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service/auth"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/workers"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

//@title URL Shortener API
//@version 1.0
//@description API Server for shortener

//@host localhost:8080
//@BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// BuildVersion = определяет версию приложения
// BuildDate = определяет дату сборки
// BuildCommit = определяет коммит сборки
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	// выводит глобальную информацию
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	// инициализируем конфиг.
	configs := NewConfigs()
	configs.Parse()

	// инициализируем logger.
	logs := logger.NewLogger(logger.WithLevel(configs.LogLevel))
	logs.Info("Start logger")

	// инициализируем хранилище.
	repo := initDB(configs.AddrDB, configs.PathFile)
	logs.Info("Connecting DB")
	defer repo.Close()

	// инициализируем сервис.
	urlService := service.NewService(
		repo,
		logs,
	)
	logs.Info("Service created")

	// инициализируем проверку авторизацию.
	serviceAuth := auth.NewServiceAuth(repo)
	authorization := middleware.NewAuthMiddleware(serviceAuth)

	// инициализируем worker.
	worker := workers.NewWorkerDeleted(urlService)

	// передаем в хенлер сервис и baseURL.
	shortHandlers := handlers.NewHandlers(urlService, configs.BaseURL, logs, worker)
	logs.Info(fmt.Sprintf("Handlers created PORT: %s", configs.AddrServer))

	// инициализировали роутер и создали Post и Get.
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)

	// Swagger route.
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Добавляем pprof маршруты вручную.
	// Ограничим доступ с помощью контроля ip.
	r.Group(func(r chi.Router) {
		r.Use(middleware.PprofMiddleware)
		r.Mount("/debug", middleware2.Profiler())
	})

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

	// слушаем выбранны порт = configs.AddrServer.
	if *configs.HTTPS {
		logs.Info("Starting HTTPS server with self-signed certificate")
		// Настройка HTTPS-сервера с самоподписанным сертификатом
		err := http.ListenAndServeTLS(configs.AddrServer, "./cert.pem", "./key.pem", r)
		if err != nil {
			logs.Error("Failed to start HTTPS server:", logger.ErrAttr(err))
		}
	} else {
		logs.Info("Starting HTTP server")
		if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
			logs.Error("Failed to start HTTP server:", logger.ErrAttr(err))
		}
	}
}
