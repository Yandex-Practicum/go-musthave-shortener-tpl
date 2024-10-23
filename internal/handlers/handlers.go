package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/errorscustom"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/middleware"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/workers"
)

type Handlers struct {
	service *service.Service
	baseURL string
	logger  *logger.Logger
	worker  workers.Worker
	//worker  *workers.WorkerDeleted
}

func NewHandlers(service *service.Service, baseURL string, sLog *logger.Logger, worker workers.Worker) *Handlers {
	return &Handlers{
		service: service,
		baseURL: baseURL,
		logger:  sLog,
		worker:  worker,
	}
}

// PostJSON godoc
// @Tags POST
// @Summary Create new short URL from JSON request
// @Description Create a short URL based on the given JSON payload
// @Accept json
// @Produce json
// @Param url body models.URL true "URL to shorten"
// @Success 201 "Created"
// @Failure 400 "Bad request"
// @Failure 404 "URL not found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal server error"
// @Router /api/shorten [post]
// PostJSON обрабатываем JSON запрос и возвращаем короткую ссылку.
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
		if errors.Is(err, errorscustom.ErrConflict) {
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

// PostURL godoc
// @Tags POST
// @Summary Create new short URL from URL
// @Description Create a short URL based on the given URL
// @Accept plain
// @Produce plain
// @Param url body string true "URL to shorten"
// @Success 201 "Created"
// @Failure 400 "Bad request"
// @Failure 404 "URL not found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal server error"
// @Router / [post]
// PostURL обрабатываем обычный запрос и возвращаем короткую ссылку.
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
		if errors.Is(err, errorscustom.ErrConflict) {
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

// PostBatchDB godoc
// @Tags POST
// @Summary Create new short URL from URL
// @Description Create a short URL based on the given URL
// @Accept json
// @Produce json
// @Param url body []models.MultipleURL true "URL to shorten"
// @Success 201 "Created"
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /api/shorten/batch [post]
// PostBatchDB записываем запрос в db.
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

// GetURL godoc
// @Tags GET
// @Summary Get short URL
// @Description Get short URL
// @Accept json
// @Produce json
// @Param id path string true "Short URL"
// @Success 307 "Temporary redirect"
// @Header 307 {string} Location "URL новой записи"
// @Failure 404 "Not found"
// @Failure 405 "Method not allowed"
// @Failure 410 "Gone"
// @Router /{id} [get]
// GetURL возвращаем информацию по короткой ссылке.
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
		if errors.Is(err, errorscustom.ErrDeletedURL) {
			h.logger.Error("error =", "GET/{id}", errorscustom.ErrDeletedURL)
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

// GetPing godoc
// @Tags GET
// @Summary Check DB connection
// @Description Check DB connection
// @Accept plain
// @Produce plain
// @Success 200 "OK"
// @Failure 500 "Internal server error"
// @Router /ping [get]
// GetPing Проверяем подключение к DB.
func (h *Handlers) GetPing(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(); err != nil {
		h.logger.Error("Error = ", logger.ErrAttr(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetUsersURLs godoc
// @Tags GET
// @Summary Get user URLs
// @Description Get user URLs
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 "OK"
// @Success 204 "No content"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/user/urls [get]
// GetUsersURLs возвращаем все сохраненные URL пользователя.
func (h *Handlers) GetUsersURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		h.logger.Error("Error = ", logger.ErrAttr(errorscustom.ErrUserIDNotContext))
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

// DeletionURLs godoc
// @Tags DELETE
// @Summary Delete user URLs
// @Description Delete user URLs
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param urls body []string true "URLs"
// @Success 202 "Accepted"
// @Failure 500 "Internal server error"
// @Router /api/user/urls [delete]
// DeletionURLs делает запрос на удаление из базы.
func (h *Handlers) DeletionURLs(w http.ResponseWriter, r *http.Request) {
	var urls []string
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&urls); err != nil {
		h.logger.Error("cannot decode request JSON body:", "error = ", err)
		http.Error(w, "cannot decode request JSON body", http.StatusInternalServerError)
		return
	}

	// получаем из контектса userID
	userID, _ := r.Context().Value(middleware.UserIDContextKey).(string)

	req := workers.DeletionRequest{
		User: userID,
		URLs: urls,
	}

	if err := h.worker.SendDeletionRequestToWorker(req); err != nil {
		h.logger.Error("error send to deletion worker request", "error = ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) ResultBody(res string) string {
	return h.baseURL + "/" + res
}
