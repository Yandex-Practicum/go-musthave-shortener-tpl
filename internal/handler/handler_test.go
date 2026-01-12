package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/klyakssa/go-musthave-shortener-tpl/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandler(t *testing.T) {
	type want struct {
		codePost    int
		codeGet     int
		url         string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "test #1",
			want: want{
				codePost:    201,
				codeGet:     307,
				url:         "https://yandex.ru",
				contentType: "text/plain",
			},
		},
		{
			name: "test #2",
			want: want{
				codePost:    201,
				codeGet:     307,
				url:         "https://practicum.yandex.ru",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.want.url))
			w := httptest.NewRecorder()
			handler.ShortenHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.codePost, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			assert.Equal(t, strconv.Itoa(len(resBody)), res.Header.Get("Content-Length"))

			require.NoError(t, err)

			request2 := httptest.NewRequest(http.MethodGet, string(resBody), nil)
			w2 := httptest.NewRecorder()
			handler.UnshortenHandler(w2, request2)

			res2 := w2.Result()

			assert.Equal(t, test.want.codeGet, res2.StatusCode)
			assert.Equal(t, test.want.url, res2.Header.Get("Location"))

			defer res2.Body.Close()
		})
	}
}
