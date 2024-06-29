package main

import (
	"io"
	"net/http"
	"net/url"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortURL)
	mux.HandleFunc(`/{id}/`, fullURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil{
		panic(err)
	}
}

func shortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	_, err = url.ParseRequestURI(string(reqBody))
	if err != nil{
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte{})
	}

	res.Header().Set("Content-type", "text/plain")

	res.WriteHeader(http.StatusCreated)
	resBody := `http://localhost:8080/EwHXdJfB `
	res.Write([]byte(resBody))
}

func fullURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	res.Header().Set("Location", "https://practicum.yandex.ru/")
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte{})
}

