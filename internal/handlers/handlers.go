package handlers

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"io"
	"log"
	"net/http"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) PostURL(w http.ResponseWriter, r *http.Request) {
	log.Println("POST URL")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
       "response": {
           "text": "Извините, я пока ничего не умею"
       },
       "version": "1.0"
   }`))
		return
	}

	encodeURL, err := h.service.SaveURL(string(body))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("URL encoded successfully")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("https://localhost:8080/" + encodeURL))
}

func (h *Handlers) GetURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GET URL")

	shortURL := r.URL.String()[1:]

	if shortURL == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	url, err := h.service.GetURL(shortURL)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(url))
	log.Printf("URL decoded successfully: %s", url)

}
