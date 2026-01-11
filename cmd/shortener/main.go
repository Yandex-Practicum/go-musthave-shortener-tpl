package main

import (
	"github.com/klyakssa/go-musthave-shortener-tpl/internal/handler"
	"github.com/klyakssa/go-musthave-shortener-tpl/internal/router"
)

func main() {
	r := router.NewMyRouter()
	r.HandleFunc("/", handler.MainHandler)
	
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
