package filestorage

import (
	"os"
	"testing"
)

// TestSaveFile_SaveURL - тестируем сохранение в файл.
func TestSaveFile_SaveURL(t *testing.T) {
	storage, err := NewSaveFile("testStorage.txt")
	if err != nil {
		t.Error("problam create test file")
	}

	shortURL := "qwert"
	originURL := "https://www.ya.ru"
	userID := "test"

	err = storage.SaveURL(shortURL, originURL, userID)

	if err != nil {
		t.Error("problam save url")
	}
	os.Remove("testStorage.txt")
}
