package main

import (
	"log"
	"net/http"

	"github.com/IgorGreusunset/shortener/cmd/config"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

//Переменные используем в качестве БД
var Storage []model.URL = []model.URL{}



func main() {
	router := chi.NewRouter()

	router.Post(`/`, handlers.PostHandler)
	router.Get(`/{id}`, handlers.GetByIDHandler)

	conf := config.ParseFlag()

	serverAdd := "http://"+conf.Serv

	log.Fatal(http.ListenAndServe(serverAdd, router))
}

