package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/klyakssa/go-musthave-shortener-tpl/internal/repository"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shrt, err := repository.Shorten(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len("http://localhost:8080/"+shrt)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shrt))
}

func UnshortenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	lng, err := repository.Unshorten(r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", lng)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
