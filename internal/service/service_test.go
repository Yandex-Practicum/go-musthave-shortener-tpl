package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("create_service", func(t *testing.T) {
		logs := logger.NewLogger(logger.WithLevel("info"))
		s := NewService(nil, logs)
		assert.NotNil(t, s)
	})
}
