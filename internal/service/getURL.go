package service

func (s *Service) GetURL(shortURL string) (string, error) {
	url, err := s.storage.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return url, nil
}
