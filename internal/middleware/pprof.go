package middleware

import "net/http"

// PprofMiddleware godoc
// @Tags MIDDLEWARE
// @Summary Pprof middleware
// @Description Pprof middleware - work only location
// @Success 200 "OK"
// @Failure 403 "Access denied"
// @Router /pprof/... [get]
// PprofMiddleware - отвечает за закрытость использования pprof.
func PprofMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pprof/" {
			if r.RemoteAddr == "127.0.0.1" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello from PProf!"))
				return
			}
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied."))
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
