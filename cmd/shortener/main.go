package main

import (
	"net/http"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/handlers"
)

//Переменные используем в качестве БД
var Storage []model.URL = []model.URL{}



func main() {
	//Storage = []model.URL{}
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.PostHandler)
	mux.HandleFunc(`/{id}`, handlers.GetByIDHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil{
		panic(err)
	}
}

