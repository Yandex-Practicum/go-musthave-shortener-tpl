package main

import (
	"fmt"

	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/db"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
)

// initDB инициализация базы.
func initDB(addrDB, pathFile string) service.Storage {

	if addrDB != "" {
		// Хранение в базе данных
		fmt.Println("Using database storage with DSN:", addrDB)
		// Инициализация базы данных и работа с ней
		repoDB, err := db.NewPstStorage(addrDB)
		if err != nil {
			fmt.Println("Fatal: ", err)
		}
		return repoDB
	} else if pathFile != "" {
		// Хранение в файле
		fmt.Println("Using database storage with file:", pathFile)
		// Инициализируем хранение в файле
		repoFile, err := filestorage.NewSaveFile(pathFile)
		if err != nil {
			fmt.Println("Fatal: ", err)
		}
		return repoFile
	} else {
		// Хранение в памяти
		fmt.Println("Using in-memory storage")
		// Инициализация хранения в памяти
		repoMap := mapstorage.NewMapURL()
		return repoMap
	}
}
