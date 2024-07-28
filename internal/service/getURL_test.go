package service

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/fileStorage"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/mapstorage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetURL(t *testing.T) {
	logs := logger.NewLogger(logger.WithLevel("info"))
	storageURL := mapstorage.NewMapURL()
	// инициализируем файл для хранения
	fileName := "./test.txt"
	defer os.Remove(fileName)

	file, err := fileStorage.NewSaveFile(fileName)
	if err != nil {
		logs.Error("Fatal", logger.ErrAttr(err))
	}
	defer file.Close()
	service := NewService(storageURL, file)
	t.Run("get URL", func(t *testing.T) {
		url, err := service.GetURL("")
		assert.NotNil(t, err)
		assert.Equal(t, "", url)

	})
}
