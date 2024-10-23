package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
)

// Service - сервис.
type Service struct {
	storage Storage
	logger  *logger.Logger
}

// NewService - конструктор сервиса.
func NewService(storage Storage, logger *logger.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}
