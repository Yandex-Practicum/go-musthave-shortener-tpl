package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/klyakssa/go-musthave-shortener-tpl/internal/repository"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		lng, err := repository.Unshorten(r.URL.Path[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", lng)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if r.Method == "POST" {
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

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(shrt)))
		w.Write([]byte("http://localhost:8080/" + shrt))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
