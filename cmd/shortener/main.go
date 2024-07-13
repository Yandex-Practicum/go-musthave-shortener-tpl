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
	log.Println("Storage created")
	urlService := service.NewService(storage)
	log.Println("Service created")
	configs := NewConfigs()
	configs.ParseFlags()
	shortHandlers := handlers.NewHandlers(urlService)
	log.Println("Handlers created")
	r := chi.NewRouter()
	r.Post("/", shortHandlers.PostURL)
	r.Get("/{url}", shortHandlers.GetURL)
	if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
		log.Fatal(err)
	}
}
