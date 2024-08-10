package service

import "github.com/kamencov/go-musthave-shortener-tpl/internal/logger"

type Service struct {
	storage Storage
	logger  *logger.Logger
}

func NewService(storage Storage, logger *logger.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}
