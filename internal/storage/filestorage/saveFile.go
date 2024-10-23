package filestorage

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

// SaveURL - функция для записи в файл.
func (s *SaveFile) SaveURL(shortURL, originalURL, userID string) error {
	Count++
	var events = []*Event{
		{
			UUID:        Count,
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		},
	}

	for _, event := range events {
		if err := s.WriteSaveModel(event); err != nil {
			return err
		}
	}
	return nil
}

// SaveSliceOfDB - функция для записи в файл.
func (s *SaveFile) SaveSliceOfDB(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return nil, nil
}
