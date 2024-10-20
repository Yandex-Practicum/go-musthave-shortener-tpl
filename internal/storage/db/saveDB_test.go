package db

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPstStorage_SaveURL(t *testing.T) {

	tests := []struct {
		name        string
		shortURL    string
		userID      string
		originalURL string
		exec        bool
		execErr     error
		expectedErr error
	}{
		{
			name:        "successful",
			shortURL:    "shortURL",
			userID:      "testID",
			originalURL: "www.test.ru",
		},
		{
			name:        "incorrect_exec",
			shortURL:    "shortURL",
			userID:      "testID",
			originalURL: "www.test.ru",
			exec:        true,
			execErr:     sql.ErrTxDone,
			expectedErr: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			mock.ExpectBegin()
			if !tt.exec {
				mock.ExpectExec("INSERT INTO urls").
					WithArgs("shortURL", "www.test.ru", "testID").
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO urls").
					WithArgs("shortURL", "www.test.ru", "testID").
					WillReturnError(tt.execErr)
				mock.ExpectRollback()
			}

			storage := &PstStorage{
				storage: db,
			}

			err = storage.SaveURL(tt.originalURL, tt.shortURL, tt.userID)

			if errors.Is(err, tt.expectedErr) {
				t.Errorf("SaveUrl = %t, want = %t", err, tt.expectedErr)
			}
		})
	}
}
