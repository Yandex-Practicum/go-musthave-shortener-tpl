package main

import (
	"log"
	"net/http"

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

	log.Fatal(http.ListenAndServe(":8080", router))
}

