package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/klyakssa/go-musthave-shortener-tpl/internal/config"
)

type MyRouter struct {
	Engine *gin.Engine
	Config *config.Config
}

func NewMyRouter(cfg *config.Config) *MyRouter {
	return &MyRouter{
		Engine: gin.Default(),
		Config: cfg,
	}
}

func (r *MyRouter) Run(addr string) error {
	return r.Engine.Run(addr)
}

func (r *MyRouter) GET(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.Engine.GET(pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}

func (r *MyRouter) POST(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.Engine.POST(pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}
