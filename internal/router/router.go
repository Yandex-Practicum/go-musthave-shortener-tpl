package router

import (
	"net/http"
)

type MyRouter struct {
	*http.ServeMux
}

func NewMyRouter() *MyRouter {
	return &MyRouter{
		ServeMux: http.NewServeMux(),
	}
}

func (r *MyRouter) Run(addr string) error {
	return http.ListenAndServe(addr, r.ServeMux)
}

func (r *MyRouter) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.ServeMux.HandleFunc(pattern, handler)
}
