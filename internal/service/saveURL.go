package service

import (
	errors2 "github.com/kamencov/go-musthave-shortener-tpl/internal/errorscustom"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/utils"
)

func (s *Service) SaveURL(url, userID string) (string, error) {

	// проверяем есть ли в базе уже данный URL
	if shortURL, err := s.storage.CheckURL(url); err != nil {
		return shortURL, errors2.ErrConflict
	}

	// создаем короткую ссылку так как не нашли в базе
	encodeURL, err := utils.EncodeURL(url)

	if err != nil {
		s.logger.Error("Error = ", logger.ErrAttr(err))
		return "", err
	}

	err = s.storage.SaveURL(encodeURL, url, userID)
	if err != nil {
		s.logger.Error("Error = ", logger.ErrAttr(err))
		return "", err
	}

	return encodeURL, nil
}
