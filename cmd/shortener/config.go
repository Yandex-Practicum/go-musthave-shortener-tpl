package main

import (
	"flag"
	"os"
)

type Configs struct {
	AddrServer   string
	BaseURL      string
	flagLogLevel string
}

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) Parse() {

	c.parseFlags()

	serverAdd := os.Getenv("SERVER_ADDRESS")
	if serverAdd != "" {
		c.AddrServer = serverAdd
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL != "" {
		c.BaseURL = baseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.flagLogLevel = envLogLevel
	}

}

func (c *Configs) parseFlags() {
	// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8080).
	flag.StringVar(&c.AddrServer, "a", ":8080", "Server address host:port")
	//Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL,
	//например http://localhost:8080/qsd54gFg).
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result net address host:port")

	flag.StringVar(&c.flagLogLevel, "l", "info", "log level")
	flag.Parse()
}
