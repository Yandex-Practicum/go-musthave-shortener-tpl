package main

import (
	"github.com/klyakssa/go-musthave-shortener-tpl/internal/config"
	"github.com/klyakssa/go-musthave-shortener-tpl/internal/handler"
	"github.com/klyakssa/go-musthave-shortener-tpl/internal/router"
)

func main() {
	config := config.InitFlagConfig()

	r := router.NewMyRouter(config)

	h := handler.NewMyHandler(config)
	r.GET("/:uuid", h.UnshortenHandler)
	r.POST("/", h.ShortenHandler)

	if err := r.Run(config.WebConfig.HostPort); err != nil {
		panic(err)
	}
}
