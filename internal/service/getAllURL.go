package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/models"

// GetAllURL возвращает все сохраненные пользователем URL-адреса
func (s *Service) GetAllURL(userID, baseURL string) ([]*models.UserURLs, error) {
	return s.storage.GetAllURL(userID, baseURL)
}
