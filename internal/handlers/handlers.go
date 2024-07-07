package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/IgorGreusunset/shortener/cmd/config"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/helpers"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(string(reqBody))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	id := helpers.Generate()

	urlToAdd := model.NewURL(id, string(reqBody))
	storage.WriteToStorage(*urlToAdd)

	res.Header().Set("Content-type", "text/plain")

	res.WriteHeader(http.StatusCreated)

	conf := config.ParseFlag()

	resBody := conf.Base + `/` + id
	res.Write([]byte(resBody))
}

func GetByIDHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}


	short := chi.URLParam(req, "id")

	fullURL := storage.ReadFromStorage(short)

	if fullURL.FullURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", fullURL.FullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte{})
}