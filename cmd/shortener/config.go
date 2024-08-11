package main

import (
	"flag"
	"fmt"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/db"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"os"
)

type Configs struct {
	AddrServer string
	BaseURL    string
	LogLevel   string
	PathDB     string
	AddrDB     string
	repository service.Storage
}

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) Parse() {
	var err error
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
		c.PathDB = pathFile
	}
	// Проверка переменной окружения DATABASE_DSN
	if envAddrDB := os.Getenv("DATABASE_DSN"); envAddrDB != "" {
		c.AddrDB = envAddrDB
	}

	if c.AddrDB != "" {
		// Хранение в базе данных
		fmt.Println("Using database storage with DSN:", c.AddrDB)
		// Инициализация базы данных и работа с ней
		c.repository, err = db.NewPstStorage(c.AddrDB)
		if err != nil {
			fmt.Println("Fatal: ", err)
		}
	} else if c.PathDB != "" {
		// Хранение в файле
		fmt.Println("Using database storage with file:", c.PathDB)
		// Инициализируем хранение в файле
		c.repository, err = filestorage.NewSaveFile(c.PathDB)
		if err != nil {
			fmt.Println("Fatal: ", err)
		}
	} else {
		// Хранение в памяти
		fmt.Println("Using in-memory storage")
		// Инициализация хранения в памяти
		c.repository = mapstorage.NewMapURL()
	}

}

func (c *Configs) parseFlags() {
	// Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8080).
	flag.StringVar(&c.AddrServer, "a", ":8080", "Server address host:port")
	//Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL,
	//например http://localhost:8080/qsd54gFg).
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Result net address host:port")
	//Флаг -f отвечает за базовый путь сохранения storage
	flag.StringVar(&c.PathDB, "f", "", "full name for file repository")
	// Флаг -l отвечает за logger
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	//Флаг -p отвечает за адрес подключения DB
	flag.StringVar(&c.AddrDB, "d", "", "address DB")
	flag.Parse()
}
