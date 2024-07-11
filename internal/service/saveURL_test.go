package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapStorage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_SaveURL(t *testing.T) {
	storageUrl := mapStorage.NewMapUrl()
	service := NewService(storageUrl)

	t.Run("save URL", func(t *testing.T) {
		_, err := service.SaveURL("")
		assert.NotNil(t, err)
		assert.Equal(t, "URL is empty", err.Error())

		_, err = service.SaveURL("http://example.com")
		assert.Nil(t, err)
	})
}
