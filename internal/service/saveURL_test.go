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

func BenchmarkService_SaveURL(b *testing.B) {
	cntl := gomock.NewController(b)
	defer cntl.Finish()
	mockStorage := mocks.NewMockStorage(cntl)
	mockStorage.EXPECT().CheckURL(gomock.Any()).Return("https://example.com", nil).AnyTimes()
	mockStorage.EXPECT().SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	service := NewService(mockStorage, logger.NewLogger(logger.WithLevel("info")))

	for i := 0; i < b.N; i++ {
		service.SaveURL("https://example.com", "")
	}
}
