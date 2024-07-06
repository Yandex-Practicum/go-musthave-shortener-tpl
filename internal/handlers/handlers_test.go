package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	model "github.com/IgorGreusunset/shortener/internal/app"
	"github.com/IgorGreusunset/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	storage.Storage = []model.URL{}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "normal test",
			body: "https://mail.ru/",
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name: "not url test",
			body: "some text not url",
			want: want{code: 400},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := strings.NewReader(test.body)
			request := httptest.NewRequest(http.MethodPost, `/`, reqBody)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostHandler)
			h(w, request)

			result := w.Result()

			defer result.Body.Close()

			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}

func TestGetByIDHandler(t *testing.T) {
	storage.Storage = []model.URL{
		{ID: "U8rtGB25", FullURL: "https://practicum.yandex.ru/"},
		{ID: "g7RETf01", FullURL: "https://mail.ru/"},
	}

	type want struct {
		code int
		location string
	}

	tests := []struct{
		name string
		request string
		want want
	}{
		{
			name: "normal 1",
			request: "/U8rtGB25",
			want: want{code: 307, location: "https://practicum.yandex.ru/"},
		},
		{
			name: "normal 2",
			request: "/g7RETf01/",
			want: want{code: 307, location: "https://mail.ru/"},
		},
		{
			name: "not found",
			request: "/yyokley",
			want: want{code: 400},
		},
	}

	for _, test := range tests{
		t.Run(test.name, func (t *testing.T)  {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetByIDHandler)
			h(w, request)

			result := w.Result()

			defer result.Body.Close()
			
			assert.Equal(t, test.want.code, result.StatusCode)
			assert.Equal(t, test.want.location, result.Header.Get("Location"))
		})
	}
}
