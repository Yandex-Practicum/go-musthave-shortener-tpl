package filestorage

import (
	"fmt"
	"os"
	"testing"
)

// TestSaveFile_SaveURL - тестируем сохранение в файл.
func TestSaveFile_SaveURL(t *testing.T) {
	storage, err := NewSaveFile("testStorage.txt")
	if err != nil {
		fmt.Errorf("problam create test file")
	}

	shortURL := "qwert"
	originURL := "https://www.ya.ru"
	userID := "test"

	err = storage.SaveURL(shortURL, originURL, userID)

	if err != nil {
		fmt.Errorf("problam save url")
	}
	os.Remove("testStorage.txt")
}
