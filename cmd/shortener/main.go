package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapStorage"
	"log"
	"net/http"
)

func main() {
	storage := mapStorage.NewMapUrl()
	urlService := service.NewService(storage)
	shortHandlers := handlers.NewHandlers(urlService)
	r := chi.NewRouter()
	r.Post("/", shortHandlers.PostURL)
	r.Get("/{url}", shortHandlers.GetURL)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
