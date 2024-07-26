package main

import (
	"log"
	"net/http"
	"os"

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

	file, err := os.OpenFile(config.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error during opening file with shorten urls: %v", err)
	}

	db := storage.NewStorage(map[string]model.URL{})

	err = db.FillFromFile(file)
	if err != nil {
		logger.Log.Infof("Error during reading from file with shorten urls: %v", err)
	}


	file.Close()
	

	//Обертки для handlers, чтобы использовать их в роутере
	PostHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.PostHandler(db, config.File, res, req)
	}

	GetHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.GetByIDHandler(db, res, req)
	}

	APIPostHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.APIPostHandler(db, config.File, res, req)
	}

	router.Use(middleware.WithLogging)
	router.Use(middleware.GzipMiddleware)

	router.Post(`/`, PostHandlerWrapper)
	router.Get(`/{id}`, GetHandlerWrapper)
	router.Post(`/api/shorten`, APIPostHandlerWrapper)

	serverAdd := config.Serv

	log.Fatal(http.ListenAndServe(serverAdd, router))
}


