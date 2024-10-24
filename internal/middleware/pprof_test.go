package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPprofMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		addr        string
		exectedCode int
	}{
		{
			name:        "successful_connection",
			path:        "/pprof/",
			addr:        "127.0.0.1",
			exectedCode: 200,
		},
		{
			name:        "foreign_access",
			path:        "/pprof/",
			addr:        "237.84.2.178",
			exectedCode: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.path, nil)
			r.RemoteAddr = tt.addr
			w := httptest.NewRecorder()

			handler := PprofMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			}))

			handler.ServeHTTP(w, r)

			if w.Code != tt.exectedCode {
				t.Errorf("Ожидали = %v, пришел %v", w.Code, tt.exectedCode)
			}
		})
	}
}
