package storage

import (
	"sync"

	model "github.com/IgorGreusunset/shortener/internal/app"
)

type Storage struct {
	db map[string]model.URL
	mu sync.RWMutex
}

//Фабричный метод создания нового экземпляра хранилища
func NewStorage(db map[string]model.URL) *Storage {
	return &Storage{db: db}
}

type Repository interface {
	Create(record model.URL)
	GetByID(id string) (model.URL, bool)
}

//Метод для создания новой записи в хранилище
func (s *Storage) Create(record model.URL) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[record.ID] = record
}

//Метода для получения записи из хранилища
func (s *Storage) GetByID(id string)(model.URL, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.db[id]
	return url, ok
}
