package filestorage

import (
	"os"
	"testing"
)

func TestSaveFile_CheckURL(t *testing.T) {
	storage, err := NewSaveFile("testStorage.txt")
	if err != nil {
		t.Error("problam create test file")
	}
	defer os.Remove("testStorage.txt")

	_, err = storage.CheckURL("qwert")
	if err != nil {
		t.Error("problam check url")
	}

}
