package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/logger"

func (s *Service) GetURL(shortURL string) (string, error) {
	url, err := s.storage.GetURL(shortURL)
	if err != nil {
		s.logger.Error("Error = ", logger.ErrAttr(err))
		return "", err
	}
	return url, nil
}
