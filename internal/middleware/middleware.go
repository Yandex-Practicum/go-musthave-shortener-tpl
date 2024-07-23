package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/IgorGreusunset/shortener/internal/compress"
	"github.com/IgorGreusunset/shortener/internal/logger"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}

		uri := req.RequestURI
		method := req.Method

		h.ServeHTTP(&lw, req)

		duration := time.Since(start)

		logger.Log.Infoln(
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
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(starusCode int) {
	r.ResponseWriter.WriteHeader(starusCode)
	r.responseData.status = starusCode
}


func GzipMiddleware(h http.Handler) http.Handler{
	comp := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		content := w.Header().Get("Content-Type")

		if content == "application/json" || content == "text/html" {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip"){
				cw := compress.NewCompressWrite(w)
				ow = cw
				defer cw.Close()			
			}
		}

		contentEncoding := r.Header.Get("Content-Encoding")

		if strings.Contains(contentEncoding, "gzip") {
			cr, err := compress.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(comp)
}