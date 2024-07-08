package handlers

import (
	. "github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapStorage"
	"io"
	"log"
	"net/http"
)

func PostURL(w http.ResponseWriter, r *http.Request) {
	log.Println("POST URL")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if string(body) != "" {
		encodURL, err := EncodeURL(string(body))

		if err != nil {
			log.Println(err)
			panic(err)
		}

		MapStorage.Storage[encodURL] = string(body)
		log.Println("URL encoded successfully")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("https://localhost:8080/" + encodURL))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GET URL")

	shortURL := r.URL.String()[1:]

	if MapStorage.Storage[shortURL] != "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(MapStorage.Storage[shortURL]))
		log.Printf("URL decoded successfully: %s", MapStorage.Storage[shortURL])
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
