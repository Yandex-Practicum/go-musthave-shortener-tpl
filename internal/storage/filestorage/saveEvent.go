package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
)

var Count int

type Event struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type SaveFile struct {
	file    *os.File
	encoder *json.Encoder
}

func NewSaveFile(filePath string) (*SaveFile, error) {
	// откройте файл и создайте для него json.Encoder
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	readFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		Count++
	}

	return &SaveFile{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (s *SaveFile) WriteSaveModel(event *Event) error {
	return s.encoder.Encode(&event)
}

func (s *SaveFile) Close() error {
	return s.file.Close()
}

type ReadFile struct {
	file    *os.File
	decoder *json.Decoder
}

func NewReadFile(filename string) (*ReadFile, error) {
	// откройте файл и создайте для него json.Decoder
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ReadFile{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil

}

func (c *ReadFile) ReadEvent() (*Event, error) {
	// добавьте вызов Decode для чтения и десериализации
	event := &Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *ReadFile) Close() error {
	return c.file.Close()
}
