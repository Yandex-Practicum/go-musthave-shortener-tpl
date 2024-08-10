package filestorage

func (s *SaveFile) SaveURL(shortURL, originalURL string) error {
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
