package main

import (
	"log"
	"net/http"
	"time"

	"github.com/IgorGreusunset/shortener/cmd/config"
	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/handlers"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)


var sugar zap.SugaredLogger

func main() {

	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar = *logger.Sugar()

	router := chi.NewRouter()

	db := storage.NewStorage(map[string]model.URL{})

	//Обертки для handlers, чтобы использовать их в роутере
	PostHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.PostHandler(db, res, req)
	}

	GetHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		handlers.GetByIDHandler(db, res, req)
	}

	router.Use(WithLogging)

	router.Post(`/`, PostHandlerWrapper)
	router.Get(`/{id}`, GetHandlerWrapper)

	config.ParseFlag()

	serverAdd := config.Serv

	log.Fatal(http.ListenAndServe(serverAdd, router))
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func (res http.ResponseWriter, req *http.Request)  {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size: 0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: res,
			responseData: responseData,
		}

		uri := req.RequestURI
		method := req.Method

		h.ServeHTTP(&lw, req)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
		
	}
	return http.HandlerFunc(logFn)
} 


type (
	responseData struct {
		status int
		size int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write (b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(starusCode int) {
	r.ResponseWriter.WriteHeader(starusCode)
	r.responseData.status = starusCode
}