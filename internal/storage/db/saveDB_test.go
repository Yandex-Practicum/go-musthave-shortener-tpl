package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/stretchr/testify/assert"
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

func TestPstStorage_SaveSliceOfDB(t *testing.T) {

	tests := []struct {
		name        string
		urls        []models.MultipleURL
		baseURL     string
		userID      string
		expectedErr error
	}{
		{
			name: "successful",
			urls: []models.MultipleURL{
				{
					CorrelationID: "1",
					OriginalURL:   "www.test.ru",
				},
			},
			baseURL: "http://localhost:8080",
			userID:  "testID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			storage := &PstStorage{
				storage: db,
			}
			mock.ExpectBegin()
			// Сначала вызывается CheckURL и возвращает ошибку, чтобы мы перешли к генерации короткой ссылки
			mock.ExpectQuery("SELECT short_url FROM urls WHERE original_url = ?").
				WithArgs(tt.urls[0].OriginalURL).
				WillReturnError(sql.ErrNoRows) // URL не найден в базе

			// Мокируем успешную генерацию короткой ссылки и сохранение в базе
			mock.ExpectExec("INSERT INTO urls").
				WithArgs("encodedURL1", tt.urls[0].OriginalURL, tt.userID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			// Мокируем завершение транзакции
			mock.ExpectRollback() // Откат транзакции для обоих вызовов SaveURL
			mock.ExpectCommit()

			_, err = storage.SaveSliceOfDB(tt.urls, tt.baseURL, tt.userID)

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
