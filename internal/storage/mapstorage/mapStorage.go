package mapstorage

import (
	"errors"
	"sync"

	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
)

// IMapStorage - интерфейс хранилища URL-адресов.
type IMapStorage interface {
	SaveURL(shortURL, url, userID string) error
	GetURL(shortURL string) (string, error)
	Close() error
}

// MapStorage - хранилище URL-адресов.
type MapStorage struct {
	storage map[string]string
	mu      sync.RWMutex
}

// NewMapURL возвращает новый хранилище URL-адресов.
func NewMapURL() *MapStorage {
	return &MapStorage{
		storage: make(map[string]string),
	}
}

// SaveURL сохраняет URL в хранилище.
func (s *MapStorage) SaveURL(shortURL, url, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if url == "" {
		return errors.New("URL is empty")
	}

	s.storage[shortURL] = url
	return nil
}

// GetURL возвращает URL из хранилища.
func (s *MapStorage) GetURL(shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.storage[shortURL]; !ok {
		return "", errors.New("URL not found")
	}
	return s.storage[shortURL], nil
}

// Close закрывает хранилище.
func (s *MapStorage) Close() error {
	return nil
}

// Ping проверяет соединение с хранилищем.
func (s *MapStorage) Ping() error {
	return nil
}

// SaveSliceOfDB сохраняет срез URL в хранилище.
func (s *MapStorage) SaveSlice(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return nil, nil
}

// GetAllURL возвращает срез URL из хранилища.
func (s *MapStorage) GetAllURL(userID, baseURL string) ([]*models.UserURLs, error) {
	return nil, errors.New("not use GetAllURL in map")
}
