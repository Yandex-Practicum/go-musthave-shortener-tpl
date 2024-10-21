package mapstorage

import (
	"errors"
	"sync"

	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
)

type MapStorage struct {
	storage map[string]string
	mu      sync.RWMutex
}

func NewMapURL() *MapStorage {
	return &MapStorage{
		storage: make(map[string]string),
	}
}

func (s *MapStorage) SaveURL(shortURL, url, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if url == "" {
		return errors.New("URL is empty")
	}

	s.storage[shortURL] = url
	return nil
}

func (s *MapStorage) GetURL(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.storage[shortURL]; !ok {
		return "", errors.New("URL not found")
	}
	return s.storage[shortURL], nil
}

func (s *MapStorage) Close() error {
	return nil
}

func (s *MapStorage) Ping() error {
	return nil
}

func (s *MapStorage) SaveSliceOfDB(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return nil, nil
}

func (s *MapStorage) GetAllURL(userID, baseURL string) ([]*models.UserURLs, error) {
	return nil, errors.New("not use GetAllURL in map")
}
