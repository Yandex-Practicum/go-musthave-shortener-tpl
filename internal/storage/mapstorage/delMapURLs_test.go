package mapstorage

import (
	"testing"
)

// TestMapStorage_DeletedURLs - тестирует удаление urls из мапы.
func TestMapStorage_DeletedURLs(t *testing.T) {
	storage := NewMapURL()
	var testURLs []string
	testURLs = append(testURLs, "www")
	err := storage.DeletedURLs(testURLs, "test")
	if err == nil {
		t.Error("no way")
	}
}
