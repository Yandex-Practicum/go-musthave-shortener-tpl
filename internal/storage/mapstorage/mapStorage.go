package mapstorage

import (
	"errors"
	"log"
	"sync"
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

func (s *MapStorage) SaveURL(shortURL, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if url == "" {
		return errors.New("URL is empty")
	}

	s.storage[shortURL] = url
	log.Printf("Save key = %s, value = %s", shortURL, url)
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
