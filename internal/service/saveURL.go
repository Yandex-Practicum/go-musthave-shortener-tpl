package service

import (
	"errors"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/utils"
	"log"
)

func (s *Service) SaveURL(url string) (string, error) {
	encodeURL, err := utils.EncodeURL(url)

	if err != nil {
		log.Println(err)
		return "", errors.New("URL is empty")
	}

	err = s.storage.SaveURL(encodeURL, url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = s.inJSON(url, encodeURL)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return encodeURL, nil
}

func (s *Service) inJSON(URL, shotrURL string) error {
	filestorage.Count++
	var events = []*filestorage.Event{
		{
			UUID:        filestorage.Count,
			ShortURL:    shotrURL,
			OriginalURL: URL,
		},
	}

	for _, event := range events {
		if err := s.dbFile.WriteSaveModel(event); err != nil {
			return err
		}
	}
	return nil
}
