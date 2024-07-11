package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapStorage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetURL(t *testing.T) {
	storageUrl := mapStorage.NewMapUrl()
	service := NewService(storageUrl)
	t.Run("get URL", func(t *testing.T) {
		url, err := service.GetURL("")
		assert.NotNil(t, err)
		assert.Equal(t, "", url)

	})
}
