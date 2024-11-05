package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	errors2 "github.com/kamencov/go-musthave-shortener-tpl/internal/errorscustom"
	"github.com/stretchr/testify/assert"
)

func TestPstStorage_CheckURL(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		execErr     error
		expectedErr error
	}{
		{
			name:        "successful",
			originalURL: "test",
			execErr:     errors2.ErrConflict,
			expectedErr: errors2.ErrConflict,
		},
		{
			name:        "no_rows",
			originalURL: "test",
			execErr:     sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			storage := &PstStorage{
				storage: db,
			}
			mock.ExpectQuery("SELECT short_url").WithArgs(tt.originalURL).WillReturnError(tt.execErr)
			_, err = storage.CheckURL(tt.originalURL)

			// Сравнение ошибок
			if err != nil && tt.expectedErr != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("Ожидали ошибку %v, пришла ошибка %v", tt.expectedErr, err)
				}
			} else if err != tt.expectedErr {
				t.Errorf("Ожидали ошибку %v, пришла ошибка %v", tt.expectedErr, err)
			}
		})
	}
}

func TestPstStorage_CheckUser(t *testing.T) {

	tests := []struct {
		name        string
		user        string
		execErr     error
		expectedErr error
	}{
		{
			name: "successful",
			user: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			storage := &PstStorage{
				storage: db,
			}
			mock.ExpectQuery("SELECT user_id").WithArgs(tt.user).WillReturnRows(sqlmock.NewRows([]string{}))
			mock.ExpectQuery("SELECT user_id").WithArgs(tt.user).WillReturnError(tt.execErr)
			err = storage.CheckUser(tt.user)

			// Сравнение ошибок
			if err != nil && tt.expectedErr != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("Ожидали ошибку %v, пришла ошибка %v", tt.expectedErr, err)
				}
			} else if err != tt.expectedErr {
				t.Errorf("Ожидали ошибку %v, пришла ошибка %v", tt.expectedErr, err)
			}
		})
	}
}
