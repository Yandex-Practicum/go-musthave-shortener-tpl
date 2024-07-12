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

func PostHandler(db storage.Repository, res http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	_, err = url.ParseRequestURI(string(reqBody))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}


	id := helpers.Generate()

	urlToAdd := model.NewURL(id, string(reqBody))
	db.Create(*urlToAdd)

	res.Header().Set("Content-type", "text/plain")

	res.WriteHeader(http.StatusCreated)


	resBody := config.Base + `/` + id
	res.Write([]byte(resBody))
}

func GetByIDHandler(db storage.Repository, res http.ResponseWriter, req *http.Request) {

	short := chi.URLParam(req, "id")

	fullURL := db.GetByID(short)

	if fullURL.FullURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", fullURL.FullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
