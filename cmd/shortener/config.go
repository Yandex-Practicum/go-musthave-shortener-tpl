// Package main предоставляет конфигурацию при запуске.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Configs структура основных зависимостей при запуске.
type Configs struct {
	AddrServer string `json:"server_address"`
	BaseURL    string `json:"base_url"`
	LogLevel   string `json:"log_level"`
	PathFile   string `json:"file_storage_path"`
	AddrDB     string `json:"database_dsn"`
	HTTPS      *bool  `json:"enable_https"`
	ConfigFile string
}

// NewConfigs конструктор конфига.
func NewConfigs() *Configs {
	return &Configs{}
}

// Parse проверка зависимостей и перезапись Configs.
func (c *Configs) Parse() {
	// Разбираем флаги, включая путь к конфигурационному файлу
	c.parseFlags()

	// Если указан JSON-файл конфигурации, загружаем его
	if c.ConfigFile != "" {
		if err := c.loadConfig(); err != nil {
			fmt.Printf("Ошибка загрузки конфигурации из файла %s: %v\n", c.ConfigFile, err)
		}
	}

	// Переопределяем параметры значениями из переменных окружения
	c.parseEnv()
}

func (c *Configs) parseEnv() {
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

	// Проверка переменной окружения ENABLE_HTTPS
	if envHTTPS := os.Getenv("ENABLE_HTTPS"); envHTTPS == "true" {
		*c.HTTPS = true
	}

}

// loadConfig загружает конфигурационный файл.
func (c *Configs) loadConfig() error {
	file, err := os.Open(c.ConfigFile)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
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

	// Флаг -s отвечает за включение HTTPS в веб версии
	c.HTTPS = flag.Bool("s", false, "Enable HTTPS")

	// Флаг -c/-config отвечает за парсинг конфигурационного JSON
	flag.StringVar(&c.ConfigFile, "c", "", "config file")
	flag.StringVar(&c.ConfigFile, "config", "", "config file")

	flag.Parse()
}
