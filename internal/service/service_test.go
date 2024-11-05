package service

import (
	"testing"

	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Run("create_service", func(t *testing.T) {
		logs := logger.NewLogger(logger.WithLevel("info"))
		s := NewService(nil, logs)
		assert.NotNil(t, s)
	})
}

func BenchmarkNewService(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewService(nil, logger.NewLogger(logger.WithLevel("info")))
	}
}
