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
	Create(record *model.URL)
	GetByID(id string) (model.URL, bool)
}

//Метод для создания новой записи в хранилище
func (s *Storage) Create(record *model.URL) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[record.ID] = *record
	record.UUID = len(s.db)
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


	for s.scan.Scan() {
		err := json.Unmarshal(s.scan.Bytes(), url) 
		if err != nil {
			return err
		}
		s.db[url.ID] = *url
	}

	return nil
}

func SaveToFile (url model.URL, file string) {
	fil, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer fil.Close()

	data, _ := json.Marshal(&url)
	data = append(data, '\n')
	fil.Write(data)
}
