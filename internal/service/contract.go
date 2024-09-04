package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

//go:generate mockgen -source=./contract.go -destination=../mocks/mock_storage.go -package=mocks
type Storage interface {
	SaveURL(shortURL, originalURL, userID string) error
	SaveSliceOfDB(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error)
	GetURL(string) (string, error)
	Close() error
	Ping() error
	CheckURL(string) (string, error)
	GetAllURL(userID, baseURL string) ([]*models.UserURLs, error)
	DeletedURLs(urls []string, userID string) error
}
