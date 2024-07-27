package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"github.com/IgorGreusunset/shortener/cmd/config"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/helpers"
	"github.com/IgorGreusunset/shortener/internal/logger"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Handler для обработки Post-запроса на запись новой URL структуры в хранилище
func PostHandler(db storage.Repository, file string, res http.ResponseWriter, req *http.Request) {

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	//Проверяем, что в теле запроса корректный URL-адрес
	_, err = url.ParseRequestURI(string(reqBody))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//Генерируем ID для короткой ссылки
	id := helpers.Generate()

	//Создаем новый экземпляр URL структуры и записываем его в хранилище
	urlToAdd := model.NewURL(id, string(reqBody))
	db.Create(urlToAdd)

	if err := storage.SaveToFile(*urlToAdd, file); err != nil {
		log.Printf("Error write to file: %v", err)
		return
	}

	//Записываем заголовок и тело ответа
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	resBody := config.Base + `/` + id
	if _, err := res.Write([]byte(resBody)); err != nil {
		log.Printf("Error writing response: %v\n", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
	}
}

// Handler для обработки Get-запроса на получение ссылки по ID
func GetByIDHandler(db storage.Repository, res http.ResponseWriter, req *http.Request) {

	//Получаем ID из запроса и ищем по нему URL структуру в хранилище
	short := chi.URLParam(req, "id")

	fullURL, ok := db.GetByID(short)

	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//Записываем заголовок ответа
	res.Header().Set("Location", fullURL.FullURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

//Handler для обработки json-запроса на создание новой ссылки
func APIPostHandler(db storage.Repository, file string, res http.ResponseWriter, req *http.Request) {

	//Получаем данные для создания URL модели из запроса
	var urlFromRequest model.APIPostRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&urlFromRequest); err != nil {
		logger.Log.Debugln("error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Проверяем корректость адреса в теле запроса
	_, err := url.ParseRequestURI(urlFromRequest.URL)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	id := helpers.Generate()

	//Создаем модель и записываем в storage
	urlToAdd := model.NewURL(id, urlFromRequest.URL)
	db.Create(urlToAdd)

	//Дублируем запись в файл
	if err := storage.SaveToFile(*urlToAdd, file); err != nil {
		log.Printf("Error writing result to file: %v", err)
	}

	//Формируем и сериализируем тело ответа
	result := config.Base + `/` + id
	resp := model.NewAPIPostResponse(result)
	response, err := json.Marshal(resp)
	if err != nil {
		logger.Log.Debugln("error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Записываем заголовок и тело ответа
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(response)
}
