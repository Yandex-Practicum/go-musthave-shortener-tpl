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

	/*_, err = url.Parse(string(reqBody))
	if err != nil {
		log.Printf("Error parsing decoded URI: %v\n", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}*/

	id := helpers.Generate()

	urlToAdd := model.NewURL(id, string(reqBody))
	db.Create(*urlToAdd)

	res.Header().Set("Content-type", "text/plain")

	res.WriteHeader(http.StatusCreated)

	//config.ParseFlag()

	resBody := config.Base + `/` + id
	res.Write([]byte(resBody))
}

func GetByIDHandler(db storage.Repository, res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}


	short := chi.URLParam(req, "id")

	fullURL := db.GetById(short)

	if fullURL.FullURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", fullURL.FullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte{})
}