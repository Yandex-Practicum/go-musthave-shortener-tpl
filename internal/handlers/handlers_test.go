package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	db := storage.NewStorage(map[string]model.URL{})

	PostHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		PostHandler(db, res, req)
	}

	handler := http.HandlerFunc(PostHandlerWrapper)

	srv := httptest.NewServer(handler)

	defer srv.Close()


	tests := []struct {
		method string
		reqBody string
		expectedCode int
		expectedContent string
	}{
		{
			method: http.MethodPost,
			reqBody: "https://mail.ru/",
			expectedCode:        201,
			expectedContent: "text/plain",
			
		},
		{
			method: http.MethodPost,
			reqBody: "some text not url",
			expectedCode: 400,
			expectedContent: "",
		},
		{
			method: http.MethodGet,
			reqBody: "https://mail.ru/",
			expectedCode: 405,
			expectedContent: "",
		},
	}
	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = test.method
			req.URL = srv.URL
			req.Body = test.reqBody

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, test.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, test.expectedContent, resp.Header().Get("Content-Type"))
		})
	}
}

func TestGetByIDHandler(t *testing.T) {
	db := storage.NewStorage(map[string]model.URL{
		"U8rtGB25": model.URL{ID: "U8rtGB25", FullURL: "https://practicum.yandex.ru/"},
		"g7RETf01": model.URL{ID: "g7RETf01", FullURL: "https://mail.ru/"},
	})

	GetHandlerWrapper := func (res http.ResponseWriter, req *http.Request)  {
		GetByIDHandler(db, res, req)
	}

	handler := http.HandlerFunc(GetHandlerWrapper)
	srv := httptest.NewServer(handler)
	defer srv.Close()


	tests := []struct{
		method string
		requestID string
		expectedCode int
		expectedLocation string
	}{
		{
			method: http.MethodGet,
			requestID: "U8rtGB25",
			expectedCode: 307, 
			expectedLocation: "https://practicum.yandex.ru/",
		},
		{
			method: http.MethodGet,
			requestID: "g7RETf01",
			expectedCode: 307, 
			expectedLocation: "https://mail.ru/",
		},
		{
			method: http.MethodGet,
			requestID: "yyokley",
			expectedCode: 400,
		},
		{
			method: http.MethodPatch,
			requestID: "yoyoyo",
			expectedCode: 405,
		},
	}

	for _, test := range tests{
		t.Run(test.method, func (t *testing.T)  {
			req := httptest.NewRequest(test.method, srv.URL+"/"+test.requestID, nil)

			cntx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, cntx))
			cntx.URLParams.Add("id", test.requestID)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetHandlerWrapper)
			h(w, req)

			res := w.Result()

			defer res.Body.Close()

			assert.Equal(t, test.expectedCode, res.StatusCode, "Response code didn't match expected")
			assert.Equal(t, test.expectedLocation, res.Header.Get("Location"), "Response Location didn't match expected")

		})
	}
}
