package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/helpers"
	"github.com/IgorGreusunset/shortener/internal/storage"
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

	resBody := `http://localhost:8080/`+id
	res.Write([]byte(resBody))
}

func GetByIDHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := req.URL.Path

	short, _ := strings.CutPrefix(string(path), "/")
	short, _ = strings.CutSuffix(short, "/")

	fullURL := storage.ReadFromStorage(short)

	res.Header().Set("location", fullURL.FullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte{})
}