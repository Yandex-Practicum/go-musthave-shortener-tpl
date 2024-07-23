package main

import (
	"log"
	"net/http"

	"github.com/IgorGreusunset/shortener/cmd/config"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/handlers"
	"github.com/IgorGreusunset/shortener/internal/logger"
	"github.com/IgorGreusunset/shortener/internal/middleware"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

//var sugar zap.SugaredLogger

func main() {

	router := chi.NewRouter()

	logger.Initialize()

	db := storage.NewStorage(map[string]model.URL{})

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

	router.Use(middleware.WithLogging)
	router.Use(middleware.GzipMiddleware)

	router.Post(`/`, PostHandlerWrapper)
	router.Get(`/{id}`, GetHandlerWrapper)
	router.Post(`/api/shorten`, APIPostHandlerWrapper)

	config.ParseFlag()

	serverAdd := config.Serv

	log.Fatal(http.ListenAndServe(serverAdd, router))
}


