package main

import (
	"flag"
	"os"
)

type Configs struct {
	AddrServer string
	BaseURL    string
	LogLevel   string
	PathDB     string
	AddrDB     string
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
		c.LogLevel = envLogLevel
	}
	if envPathDB := os.Getenv("FILE_STORAGE_PATH"); envPathDB != "" {
		c.PathDB = envPathDB
	}

	if envAddrDB := os.Getenv("DATABASE_DSN"); envAddrDB != "" {
		c.AddrDB = envAddrDB
	}
}

func (c *Configs) parseFlags() {
	// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8080).
	flag.StringVar(&c.AddrServer, "a", ":8080", "Server address host:port")
	//Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL,
	//например http://localhost:8080/qsd54gFg).
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result net address host:port")
	//Флаг -f отвечает за базовый путь сохранения storage
	flag.StringVar(&c.PathDB, "f", "./exampl.txt", "full name for file repository")

	flag.StringVar(&c.LogLevel, "l", "info", "log level")

	//Флаг -p отвечает за адрес подключения DB
	flag.StringVar(&c.AddrDB, "d", "", "address DB")
	flag.Parse()
}
