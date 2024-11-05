package filestorage

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

// SaveURL - функция для записи в файл.
func (s *SaveFile) SaveURL(shortURL, originalURL, userID string) error {
	Count++
	event := &Event{
		UUID:        Count,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	// Записываем событие напрямую, избегая создания массива.
	if err := s.WriteSaveModel(event); err != nil {
		return err
	}
	return nil
}

// SaveSliceOfDB - функция для записи в файл.
func (s *SaveFile) SaveSlice(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return nil, nil
}
