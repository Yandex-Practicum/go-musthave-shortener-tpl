package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/handlers"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", handlers.PostURL)
	r.Get("/{url}", handlers.GetURL)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
