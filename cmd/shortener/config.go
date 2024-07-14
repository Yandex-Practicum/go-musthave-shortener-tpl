package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Configs struct {
	AddrServer string
	BaseURL    string
	//AddrResult string
}

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) Parse() {
	//Если указана переменная окружения, то используется она.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	c.parseFlags()

	serverAdd := os.Getenv("SERVER_ADDRESS")
	if serverAdd != "" {
		c.AddrServer = serverAdd
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL != "" {
		c.BaseURL = baseURL
	}

}

func (c *Configs) parseFlags() {
	// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
	flag.StringVar(&c.AddrServer, "a", ":8080", "Server address host:port")
	//Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL,
	//например http://localhost:8000/qsd54gFg).
	flag.StringVar(&c.BaseURL, "b", "https://localhost:8000", "Result net address host:port")
	flag.Parse()
}
