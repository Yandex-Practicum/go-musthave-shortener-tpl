package main

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"strconv"
)

//Переменные используем в качестве БД
var index []string
var db map[string]string

//generateShort генерирует строку, которая будет использоваться длясокращения URL
func generateShort() string{
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
    "abcdefghijklmnopqrstuvwxyz" +
    "0123456789~=+%^*/()[]{}/!@#$?|")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
    	b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortURL)
	mux.HandleFunc(`/{id}/`, fullURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil{
		panic(err)
	}
}

//shortURl обрабатывает запрос для генерации короткого URL
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

	//генерируем строку для URL
	short := generateShort()

	//записываем в "БД"
	index = append(index, short)
	db[short] = string(reqBody)
	resBody := `http://localhost:8080/` + short
	res.Write([]byte(resBody))
}

func fullURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Получаем путь из запроса, достаем id и переводим в int
	path := req.URL.Path
	
	num := strings.Split(path, "/")
	id, err := strconv.Atoi(num[1])
	if err != nil{
		res.WriteHeader(http.StatusBadRequest)
	}

	//проверка, что есть короткая ссылка с таким id
	if id < 0 || id >= len(index) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//получаем из БД сначала короткую ссылку по индексу, а затем полную
	short := index[id]
	full, ok := db[short]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", full)
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Write([]byte{})
}

