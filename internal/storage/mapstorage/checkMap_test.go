package mapstorage

import (
	"testing"
)

// TestMapStorage_CheckURL - проверяет есть ли url в мапе уже.
func TestMapStorage_CheckURL(t *testing.T) {
	storage := NewMapURL()
	url := "www"
	_, err := storage.CheckURL(url)
	if err != nil {
		t.Error("no way")
	}

}
