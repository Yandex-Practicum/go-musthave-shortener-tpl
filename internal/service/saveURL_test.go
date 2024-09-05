package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestService_SaveURL(t *testing.T) {
	logs := logger.NewLogger(logger.WithLevel("info"))
	storageURL := mapstorage.NewMapURL()
	// инициализируем файл для хранения
	fileName := "./test.txt"
	defer os.Remove(fileName)

	file, err := filestorage.NewSaveFile(fileName)
	if err != nil {
		logs.Error("Fatal", logger.ErrAttr(err))
	}
	defer file.Close()
	service := NewService(storageURL, logs)

	t.Run("save_URL", func(t *testing.T) {
		_, err := service.SaveURL("", "")
		assert.NotNil(t, err)
		assert.Equal(t, "URL is empty", err.Error())

		_, err = service.SaveURL("http://example.com", "")
		assert.Nil(t, err)
	})
}
