package main

import (
	"flag"
	"os"
)

// Configs структура основных зависимостей при запуске.
type Configs struct {
	AddrServer string
	BaseURL    string
	LogLevel   string
	PathFile   string
	AddrDB     string
}

// NewConfigs конструктор конфига.
func NewConfigs() *Configs {
	return &Configs{}
}

// Parse проверка зависимостей и перезапись Configs.
func (c *Configs) Parse() {
	c.parseFlags()

	// Проверка переменной окружения SERVER_ADDRESS
	serverAdd := os.Getenv("SERVER_ADDRESS")
	if serverAdd != "" {
		c.AddrServer = serverAdd
	}
	// Проверка переменной окружения BASE_URL
	baseURL := os.Getenv("BASE_URL")
	if baseURL != "" {
		c.BaseURL = baseURL
	}
	// Проверка переменной окружения LOG_LEVEL
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		c.LogLevel = envLogLevel
	}
	// Проверка переменной окружения FILE_STORAGE_PATH
	if pathFile := os.Getenv("FILE_STORAGE_PATH"); pathFile != "" {
		c.PathFile = pathFile
	}
	// Проверка переменной окружения DATABASE_DSN
	if envAddrDB := os.Getenv("DATABASE_DSN"); envAddrDB != "" {
		c.AddrDB = envAddrDB
	}

}

// parseFlags флаги которые можно задать при запуске.
func (c *Configs) parseFlags() {

	// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8080).
	flag.StringVar(&c.AddrServer, "a", ":8080", "Server address host:port")

	// Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL,
	// например http://localhost:8080/qsd54gFg).
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result net address host:port")

	// Флаг -f отвечает за базовый путь сохранения storage
	flag.StringVar(&c.PathFile, "f", "", "full name for file repository")

	// Флаг -l отвечает за logger
	flag.StringVar(&c.LogLevel, "l", "info", "log level")

	// Флаг -p отвечает за адрес подключения DB
	flag.StringVar(&c.AddrDB, "d", "", "address DB")
	flag.Parse()
}
