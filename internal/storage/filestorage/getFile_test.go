package filestorage

import (
	"errors"
	"os"
	"testing"
)

func TestSaveFile_GetURL(t *testing.T) {
	for _, tt := range []struct {
		name         string
		shortURLSave string
		shortURLGate string
		expectedErr  error
	}{
		{
			name:         "successful",
			shortURLSave: "qwert",
			shortURLGate: "qwert",
			expectedErr:  nil,
		},
		{
			name:         "no_rows",
			shortURLSave: "qwert",
			shortURLGate: "not_found",
			expectedErr:  ErrShortURLNoFound,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewSaveFile("testStorage_" + tt.name + ".txt")
			if err != nil {
				t.Fatalf("ошибка создания тестового файла: %v", err)
			}
			defer os.Remove("testStorage_" + tt.name + ".txt")

			if tt.shortURLSave != "" {
				err = storage.SaveURL(tt.shortURLSave, "www.test.ru", "test")
				if err != nil {
					t.Fatalf("ошибка при сохранении URL: %v", err)
				}
			}

			_, err = storage.GetURL(tt.shortURLGate)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("ожидали ошибку %v, получили %v", tt.expectedErr, err)
			}
		})
	}
}

func TestSaveFile_GetAllURL(t *testing.T) {
	cases := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "successful",
			expectedErr: ErrNoUseInFile,
		},
	}

	storage, err := NewSaveFile("testStorage.txt")
	if err != nil {
		t.Error("problam create test file")
	}

	defer os.Remove("testStorage.txt")

	_, err = storage.GetAllURL("test", "http://localhost:8080")

	if !errors.Is(err, cases[0].expectedErr) {
		t.Errorf("ожидали ошибку %v, получили %v", cases[0].expectedErr, err)
	}

}
