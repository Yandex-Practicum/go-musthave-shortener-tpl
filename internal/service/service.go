package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/fileStorage"
)

type Service struct {
	storage Storage
	dbFile  *fileStorage.SaveFile
}

func NewService(storage Storage, dbFile *fileStorage.SaveFile) *Service {
	return &Service{
		storage: storage,
		dbFile:  dbFile,
	}
}
