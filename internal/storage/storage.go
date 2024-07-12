package storage

import (
	"sync"

	model "github.com/IgorGreusunset/shortener/internal/app"
)

type Storage struct {
	db map[string]model.URL
	mu sync.RWMutex
}

func NewStorage(db map[string]model.URL) *Storage {
	return &Storage{db: db}
}

type Repository interface {
	Create(record model.URL)
	GetByID(id string) model.URL
}

func (s *Storage) Create(record model.URL) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[record.ID] = record
}

func (s *Storage) GetByID(id string) model.URL {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, exists := s.db[id]
	if !exists {
		return model.URL{}
	}

	return url
}
