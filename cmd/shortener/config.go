package main

import (
	"flag"
)

type Configs struct {
	AddrServer string
	BaseURL    string
	//AddrResult string
}

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) ParseFlags() {
	// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
	flag.StringVar(&c.AddrServer, "a", "localhost:8080", "address server")

	//Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL,
	//например http://localhost:8000/qsd54gFg).
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8000/qsd54gFg", "address prefix")
	flag.Parse()

}
