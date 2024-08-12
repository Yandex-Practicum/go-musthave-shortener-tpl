package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
)

type Service struct {
	storage Storage
	dbFile  *filestorage.SaveFile
}

func NewService(storage Storage, dbFile *filestorage.SaveFile) *Service {
	return &Service{
		storage: storage,
		dbFile:  dbFile,
	}
}
