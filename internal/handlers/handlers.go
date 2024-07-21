package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"io"
	"net/http"
)

type Handlers struct {
	service *service.Service
	baseURL string
	logger  *logger.Logger
}

func NewHandlers(service *service.Service, baseURL string, sLog *logger.Logger) *Handlers {
	return &Handlers{
		service: service,
		baseURL: baseURL,
		logger:  sLog,
	}
}

func (h *Handlers) PostURL(w http.ResponseWriter, r *http.Request) {

	// читаем запрос из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error bad request = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// проверяем на пустой body
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

	// создаем короткую ссылку
	encodeURL, err := h.service.SaveURL(string(body))
	if err != nil {
		h.logger.Error("Error internal server = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// записываем статус, короткую сссылку
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.baseURL + "/" + encodeURL))
}

func (h *Handlers) GetURL(w http.ResponseWriter, r *http.Request) {

	shortURL := chi.URLParam(r, "id")

	if shortURL == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	url, err := h.service.GetURL(shortURL)
	if err != nil {
		//h.logger.Error("Error = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
