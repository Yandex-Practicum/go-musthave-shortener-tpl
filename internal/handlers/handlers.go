package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	errors2 "github.com/kamencov/go-musthave-shortener-tpl/internal/errors"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/templates"
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

	// получаем из контектса userID
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		h.logger.Info("Error = not userID")
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

	//

	// создаем короткую ссылку
	encodeURL, err := h.service.SaveURL(url.URL, userID)
	if err != nil {
		if errors.Is(err, errors2.ErrConflict) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.ResultURL{URL: h.ResultBody(encodeURL)})
			return
		}
		h.logger.Error("Error internal server = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// записываем заголовок, статус и короткую ссылку
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	result := models.ResultURL{
		URL: h.ResultBody(encodeURL),
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
	// получаем из контектса userID
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		h.logger.Info("Error = not userID")
	}

	// создаем короткую ссылку
	encodeURL, err := h.service.SaveURL(string(body), userID)
	if err != nil {
		if errors.Is(err, errors2.ErrConflict) {
			h.logger.Info("Conflict error: ", logger.ErrAttr(err))
			// записываем заголовок, статус и короткую ссылку
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(h.ResultBody(encodeURL)))
			return
		}
		h.logger.Error("Error internal server = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// записываем заголовок, статус и короткую ссылку
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.ResultBody(encodeURL)))
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

	// получаем из контектса userID
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		h.logger.Info("Error = not userID")
	}

	// Записываем в пустую структуру полученный запрос
	err = json.Unmarshal(body, &multipleURL)
	if err != nil {
		h.logger.Debug("cannot decode request JSON body", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultMultipleURL, err := h.service.SaveSliceOfDB(multipleURL, h.baseURL, userID)
	if err != nil {
		h.logger.Error("Error shorten URL = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
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
		if errors.Is(err, errors2.ErrDeletedURL) {
			h.logger.Error("error =", "GET/{id}", errors2.ErrDeletedURL)
			w.WriteHeader(http.StatusGone)
			return
		}
		h.logger.Error("GET/{id} =", logger.ErrAttr(err))
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetUsersURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		h.logger.Error("Error = ", logger.ErrAttr(errors2.ErrUserIDNotContext))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Получаем список URL пользователя
	listURLs, err := h.service.GetAllURL(userID, h.baseURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(listURLs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(listURLs); err != nil {
		h.logger.Error(`"error": "failed to marshal response", "details": `, logger.ErrAttr(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var sliceURLs []string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error bad request = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// получаем из контектса userID
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		h.logger.Info("Error = not userID")
		http.Error(w, "not userID", http.StatusUnauthorized)
		return
	}

	if err := json.Unmarshal(body, &sliceURLs); err != nil {
		h.logger.Debug("cannot decode request JSON body", logger.ErrAttr(err))
		http.Error(w, "cannot decode request", http.StatusBadRequest)
		return
	}

	// сигнальный канал для завершения горутин
	doneCh := make(chan struct{})
	defer close(doneCh)

	// объединяем все в один канал
	finalCH := templates.FanIn(doneCh, sliceURLs)

	// меняем статус флага в столбце is_deleted
	err = h.service.DeletedURLs(doneCh, finalCH, userID)
	if err != nil {
		h.logger.Error("Error deleted urls", logger.ErrAttr(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) ResultBody(res string) string {
	return h.baseURL + "/" + res
}
