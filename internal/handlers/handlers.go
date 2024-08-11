package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
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

// PostJSON обрабатываем JSON запрос и возвращаем короткую ссылку
func (h *Handlers) PostJSON(w http.ResponseWriter, r *http.Request) {

	// создаем структуру для сохранения URL
	url := models.URL{
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{
       "response": {
           "text": "Извините, я пока ничего не умею"
       },
       "version": "1.0"
   }`))
		return
	}

	// Записываем в пустую структуру полученный запрос
	err = json.Unmarshal(body, &url)
	if err != nil {
		h.logger.Debug("cannot decode request JSON body", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// создаем короткую ссылку
	encodeURL, err := h.service.SaveURL(url.URL)
	if err != nil {
		h.logger.Error("Error internal server = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// создаем корректную ссылку
	resultEncodingURL := h.baseURL + "/" + encodeURL

	// записываем заголовок, статус и короткую ссылку
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	result := models.ResultURL{
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

// PostURL обрабатываем обычный запрос и возвращаем короткую ссылку
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
		w.WriteHeader(http.StatusNotFound)
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

	// создаем корректную ссылку
	resultEncodingURL := h.baseURL + "/" + encodeURL

	// записываем заголовок, статус и короткую ссылку
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resultEncodingURL))
}

// PostBatchDB записываем запрос в db
func (h *Handlers) PostBatchDB(w http.ResponseWriter, r *http.Request) {
	var multipleURL []models.MultipleURL

	// читаем запрос из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error bad request = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// проверяем на пустой body
	if string(body) == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{
       "response": {
           "text": "Извините, я пока ничего не умею"
       },
       "version": "1.0"
   }`))
		return
	}

	// Записываем в пустую структуру полученный запрос
	err = json.Unmarshal(body, &multipleURL)
	if err != nil {
		h.logger.Debug("cannot decode request JSON body", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resultMultipleURL []models.ResultMultipleURL

	// создаем короткую ссылку и записываем в resultMultipleURL
	for _, req := range multipleURL {
		shortURL, err := h.service.SaveURL(req.OriginalURL)
		if err != nil {
			h.logger.Error("Error shorten URL = ", logger.ErrAttr(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resultMultipleURL = append(resultMultipleURL, models.ResultMultipleURL{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	jsonResponse, err := json.Marshal(resultMultipleURL)
	if err != nil {
		h.logger.Error("Error marshal JSON response = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

// GetURL возвращаем информацию по коротокой ссылке
func (h *Handlers) GetURL(w http.ResponseWriter, r *http.Request) {

	// читаем запрос по ключу
	shortURL := chi.URLParam(r, "id")

	//проверяем на пустой запрос
	if shortURL == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//ищем в мапе сохраненный url
	url, err := h.service.GetURL(shortURL)
	if err != nil {
		h.logger.Error("Error = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// записываем заголовок и статус
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

// GetPing Проверяем подключение к DB
func (h *Handlers) GetPing(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(); err != nil {
		h.logger.Error("Error = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
