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
		name string
		method string
		reqBody string
		expectedCode int
		expectedContent string
	}{
		{
			name: "normal case",
			method: http.MethodPost,
			reqBody: "https://mail.ru/",
			expectedCode:        http.StatusCreated,
			expectedContent: "text/plain",
			
		},
		{
			name: "not url case",
			method: http.MethodPost,
			reqBody: "some text not url",
			expectedCode: http.StatusBadRequest,
			expectedContent: "",
		},
		{
			name: "get case",
			method: http.MethodGet,
			reqBody: "https://mail.ru/",
			expectedCode: http.StatusMethodNotAllowed,
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

			if err != nil {
				t.Errorf("error making HTTP request: %v", err)
			}

			if resp.StatusCode() != test.expectedCode {
				t.Errorf("Response code didn't match expected: got %d want %d", resp.StatusCode(), test.expectedCode)
			}

			if resp.Header().Get("Content-Type") != test.expectedContent {
				t.Errorf("Response content-type didn't match expected: got %v want %v", resp.Header().Get("Content-Type"), test.expectedContent)
			}
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
		name string
		method string
		requestID string
		expectedCode int
		expectedLocation string
	}{
		{
			name: "normal practicum",
			method: http.MethodGet,
			requestID: "U8rtGB25",
			expectedCode: http.StatusTemporaryRedirect, 
			expectedLocation: "https://practicum.yandex.ru/",
		},
		{
			name: "normal mail",
			method: http.MethodGet,
			requestID: "g7RETf01",
			expectedCode: http.StatusTemporaryRedirect, 
			expectedLocation: "https://mail.ru/",
		},
		{
			name: "id not in storage",
			method: http.MethodGet,
			requestID: "yyokley",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not get method",
			method: http.MethodPatch,
			requestID: "yoyoyo",
			expectedCode: http.StatusMethodNotAllowed,
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

			if res.StatusCode != test.expectedCode {
				t.Errorf("Response code didn't match expected: got %d want %d", res.StatusCode, test.expectedCode)
			}

			if res.Header.Get("Location") != test.expectedLocation {
				t.Errorf("Response Location didn't match expected: got %v want %v", res.Header.Get("Location"), test.expectedLocation)
			}
		})
	}
}
