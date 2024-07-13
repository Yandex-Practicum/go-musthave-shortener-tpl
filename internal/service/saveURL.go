package service

import (
	"errors"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/utils"
	"log"
	"strings"
)

func (s *Service) SaveURL(url string) (string, error) {
	encodeURL, err := utils.EncodeURL(url)

	if err != nil {
		log.Println(err)
		return "", errors.New("URL is empty")
	}

	bodySplit := strings.Split(url, "\n")
	trimmedURL := bodySplit[len(bodySplit)-1]
	err = s.storage.SaveURL(encodeURL, trimmedURL)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("Save URL complete")

	return encodeURL, nil
}
