package mapstorage

import (
	"fmt"
	"testing"
)

// TestMapStorage_DeletedURLs - тестирует удаление urls из мапы.
func TestMapStorage_DeletedURLs(t *testing.T) {
	storage := NewMapURL()
	var testURLs []string
	testURLs = append(testURLs, "www")
	err := storage.DeletedURLs(testURLs, "test")
	if err == nil {
		fmt.Errorf("no way")
	}
}
