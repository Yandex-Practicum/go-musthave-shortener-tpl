package mapstorage

import (
	"fmt"
	"testing"
)

// TestMapStorage_CheckURL - проверяет есть ли url в мапе уже.
func TestMapStorage_CheckURL(t *testing.T) {
	storage := NewMapURL()
	url := "www"
	_, err := storage.CheckURL(url)
	if err != nil {
		fmt.Errorf("no way")
	}

}
