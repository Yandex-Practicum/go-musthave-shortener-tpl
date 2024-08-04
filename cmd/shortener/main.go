package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
    _ "github.com/jackc/pgx/v5/stdlib"
	"github.com/IgorGreusunset/shortener/cmd/config"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/handlers"
	"github.com/IgorGreusunset/shortener/internal/logger"
	"github.com/IgorGreusunset/shortener/internal/middleware"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)


func main() {

	config.ParseFlag()

	router := chi.NewRouter()

	logger.Initialize()


	//Открываем файл-хранилище
	file, err := os.OpenFile(config.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error during opening file with shorten urls: %v", err)
	}

	//Создаем новое хранилище
	db := storage.NewStorage(map[string]model.URL{})
	db.SetFile(file)

	var database *sql.DB

	if config.DataBase != "" {
		database, err = sql.Open("pgx", config.DataBase)
		if err != nil {
			log.Fatalf("Error during database connection: %v", err)
		}
		defer database.Close()
	}

	//Наполняем хранилище данными из файла
	err = db.FillFromFile(file)
	if err != nil {
		logger.Log.Infof("Error during reading from file with shorten urls: %v", err)
	}


	file.Close()
	

	//Обертки для handlers, чтобы использовать их в роутере
	PostHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.PostHandler(db, res, req)
	}

	GetHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.GetByIDHandler(db, res, req)
	}

	APIPostHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.APIPostHandler(db, res, req)
	}

	PingHandlerWrapper := func(res http.ResponseWriter, req *http.Request) {
		handlers.PingHandler(database, res, req)
	}

	//Подключаем middlewares
	router.Use(middleware.WithLogging)
	router.Use(middleware.GzipMiddleware)

	router.Post(`/`, PostHandlerWrapper)
	router.Get(`/{id}`, GetHandlerWrapper)
	router.Post(`/api/shorten`, APIPostHandlerWrapper)
	router.Get(`/ping`, PingHandlerWrapper)

	serverAdd := config.Serv

	log.Fatal(http.ListenAndServe(serverAdd, router))
}


