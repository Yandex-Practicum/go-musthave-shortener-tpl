package filestorage

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

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

func (s *SaveFile) SaveSliceOfDB(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return nil, nil
}
