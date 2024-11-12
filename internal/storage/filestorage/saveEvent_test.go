package filestorage

import (
	"os"
	"testing"
)

func TestNewSaveFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name     string
		nameFile string
		wantErr  bool
	}{
		{
			name:     "successful_create",
			nameFile: "test.txt",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSaveFile(tt.nameFile)
			defer os.Remove(tt.nameFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSaveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}

	t.Run("no_open_file", func(t *testing.T) {
		_, err := NewSaveFile("")
		if err == nil {
			t.Errorf("NewSaveFile() error = %v, wantErr %v", err, true)
			return
		}
	})
}
