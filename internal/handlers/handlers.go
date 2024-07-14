package handlers

import (
	"github.com/go-chi/chi/v5"
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

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + encodeURL))
}

func (h *Handlers) GetURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GET URL")
	shortURL := chi.URLParam(r, "id")
	log.Println(shortURL)
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
	log.Println(url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//w.Write([]byte(url))
	log.Printf("URL decoded successfully: %s", url)

}
