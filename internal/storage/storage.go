package storage

import model "github.com/IgorGreusunset/shortener/internal/app"

type Storage struct {
	db map[string]model.URL
}

func NewStorage(db map[string]model.URL) *Storage {
	return &Storage{db: db}
}

type Repository interface {
	Create(record model.URL)
	GetByID(id string) model.URL
}

func (s *Storage) Create(record model.URL) {
	s.db[record.ID] = record
}

func (s *Storage) GetByID(id string) model.URL {
	url, exists := s.db[id]
	if !exists {
		return model.URL{}
	}

	return url
}
