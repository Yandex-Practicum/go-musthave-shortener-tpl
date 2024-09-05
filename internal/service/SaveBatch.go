package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

func (s *Service) SaveSliceOfDB(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return s.storage.SaveSliceOfDB(urls, baseURL, userID)
}
