package db

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPstStorage_DeletedURLs(t *testing.T) {
	tests := []struct {
		name        string
		urls        []string
		userID      string
		expectedErr error
		exec        bool
	}{
		{
			name: "successful",
			urls: []string{
				"www", "ttt",
			},
			userID: "test",
		},
		{
			name: "incorrect_urls",
		},
		{
			name: "incorrect_exec",
			urls: []string{
				"www", "ttt",
			},
			exec:        true,
			expectedErr: fmt.Errorf("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock базы данных
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Инициализируем PstStorage с mock-базой
			storage := &PstStorage{
				storage: db,
			}

			if !tt.exec {
				// Мокируем ожидаемое поведение базы данных
				mock.ExpectExec("UPDATE urls").
					WithArgs(tt.userID, pq.Array(tt.urls)).WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				// Мокируем ожидаемое поведение базы данных
				mock.ExpectExec("UPDATE urls").
					WithArgs(tt.userID, pq.Array(tt.urls)).WillReturnError(fmt.Errorf("some error"))

			}
			// Вызываем тестируемую функцию
			err = storage.DeletedURLs(tt.urls, tt.userID)

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
