package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

// SaveSliceOfDB сохраняет массив коротких ссылок в базу данных
func (s *Service) SaveSliceOfDB(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	return s.storage.SaveSlice(urls, baseURL, userID)
}
