package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"io"
	"log"
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

func (h *Handlers) PostJSON(w http.ResponseWriter, r *http.Request) {
	url := struct {
		URL string `json:"url"`
	}{
		URL: "",
	}

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

	err = json.Unmarshal(body, &url)
	if err != nil {
		h.logger.Debug("cannot decode request JSON body", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encodeURL, err := h.service.SaveURL(url.URL)
	if err != nil {
		h.logger.Error("Error internal server = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultEncodingURL := h.baseURL + "/" + encodeURL

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	result := struct {
		URL string `json:"result"`
	}{
		URL: resultEncodingURL,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		h.logger.Error("error encode to json", //nolint:contextcheck // false positive
			logger.ErrAttr(err))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
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
		h.logger.Error("Error = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println(url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (h *Handlers) resultBody(body string, w http.ResponseWriter) {
	// создаем короткую ссылку
	encodeURL, err := h.service.SaveURL(body)
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
