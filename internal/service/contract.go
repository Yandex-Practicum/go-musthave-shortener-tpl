package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

//go:generate mockgen -source=./contract.go -destination=../mocks/mock_storage.go -package=mocks
type Storage interface {
	SaveURL(string, string) error
	SaveSliceOfDB(urls []models.MultipleURL, baseURL string) ([]models.ResultMultipleURL, error)
	GetURL(string) (string, error)
	Close() error
	Ping() error
	CheckURL(string) (string, error)
}
