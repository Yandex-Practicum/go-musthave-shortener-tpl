package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/utils"
)

func (s *Service) SaveURL(url string) (string, error) {
	encodeURL, err := utils.EncodeURL(url)

	if err != nil {
		s.logger.Error("Error = ", logger.ErrAttr(err))
		return "", err
	}

	err = s.storage.SaveURL(encodeURL, url)
	if err != nil {
		s.logger.Error("Error = ", logger.ErrAttr(err))
		return "", err
	}

	return encodeURL, nil
}
