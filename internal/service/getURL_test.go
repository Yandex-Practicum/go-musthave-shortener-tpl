package service

import (
	"github.com/golang/mock/gomock"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/logger"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/mocks"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/storage/filestorage"
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

	file, err := filestorage.NewSaveFile(fileName)
	if err != nil {
		logs.Error("Fatal", logger.ErrAttr(err))
	}
	defer file.Close()
	service := NewService(storageURL, logs)
	t.Run("get_URL", func(t *testing.T) {
		url, err := service.GetURL("")
		assert.NotNil(t, err)
		assert.Equal(t, "", url)

	})
}

func BenchmarkService_GetURL(b *testing.B) {

	cntl := gomock.NewController(b)
	defer cntl.Finish()
	mockStorage := mocks.NewMockStorage(cntl)
	mockStorage.EXPECT().GetURL(gomock.Any()).Return("https://example.com", nil).AnyTimes()

	service := NewService(mockStorage, logger.NewLogger(logger.WithLevel("info")))

	for i := 0; i < b.N; i++ {
		service.GetURL("")
	}
}
