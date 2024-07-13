package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"log"
	"net/http"
)

func main() {

	storage := mapstorage.NewMapURL()
	log.Println("Storage created")
	urlService := service.NewService(storage)
	log.Println("Service created")
	configs := NewConfigs()
	configs.Parse()
	shortHandlers := handlers.NewHandlers(urlService)
	log.Printf("Handlers created PORT: %s", configs.AddrServer)
	r := chi.NewRouter()
	r.Post("/", shortHandlers.PostURL)
	r.Get("/{id}", shortHandlers.GetURL)

	if err := http.ListenAndServe(configs.AddrServer, r); err != nil {
		log.Fatal(err)
	}
}
