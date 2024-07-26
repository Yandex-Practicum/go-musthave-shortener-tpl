package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	model "github.com/IgorGreusunset/shortener/internal/app"
)

type Storage struct {
	db map[string]model.URL
	mu sync.RWMutex
	w *bufio.Writer
	scan *bufio.Scanner
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
	record.UUID = len(s.db)

	data, _ := json.Marshal(record)
	s.w.Write(data)
	s.w.WriteByte('\n')
	s.w.Flush()	
}

//Метода для получения записи из хранилища
func (s *Storage) GetByID(id string)(model.URL, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.db[id]
	return url, ok
}


func (s *Storage) FillFromFile(file *os.File) error {
	url := &model.URL{}
	s.scan = bufio.NewScanner(file)
	s.w = bufio.NewWriter(file)

	for s.scan.Scan() {
		err := json.Unmarshal(s.scan.Bytes(), url) 
		if err != nil {
			return err
		}
		s.db[url.ID] = *url
	}

	return nil
}

/*func (s *Storage) SaveToFile(file *os.File) error {
	file.Seek(0, 0)
	s.w = bufio.NewWriter(file)

	for _, url := range s.db {
		data, err := json.Marshal(&url)
		if err != nil {
			return err
		}
		if _, err := s.w.Write(data); err != nil {
			return err
		}
		if err := s.w.WriteByte('\n'); err != nil {
			return err
		}
	}

	s.w.Flush()
	file.Close()
	return nil
}*/
