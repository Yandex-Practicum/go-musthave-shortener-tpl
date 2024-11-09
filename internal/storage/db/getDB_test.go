package db

import (
	"database/sql"
	errors2 "github.com/kamencov/go-musthave-shortener-tpl/internal/errorscustom"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestPstStorage_GetURL(t *testing.T) {
	// Создаем mock SQL соединение
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Инициализируем PstStorage с mock-базой данных
	pstStorage := &PstStorage{storage: db}

	tests := []struct {
		name         string
		shortURL     string
		expectedURL  string
		expectedErr  error
		mockBehavior func()
	}{
		{
			name:        "successful",
			shortURL:    "qwerty",
			expectedURL: "http://original-url.com",
			expectedErr: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"original_url", "is_deleted"}).
					AddRow("http://original-url.com", false)
				mock.ExpectQuery("SELECT original_url, is_deleted FROM urls WHERE short_url = \\$1").
					WithArgs("qwerty").
					WillReturnRows(rows)
			},
		},
		{
			name:        "url deleted",
			shortURL:    "qwerty",
			expectedURL: "",
			expectedErr: errors2.ErrDeletedURL,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"original_url", "is_deleted"}).
					AddRow("http://original-url.com", true)
				mock.ExpectQuery("SELECT original_url, is_deleted FROM urls WHERE short_url = \\$1").
					WithArgs("qwerty").
					WillReturnRows(rows)
			},
		},
		{
			name:        "url not found",
			shortURL:    "notfound",
			expectedURL: "",
			expectedErr: sql.ErrNoRows,
			mockBehavior: func() {
				mock.ExpectQuery("SELECT original_url, is_deleted FROM urls WHERE short_url = \\$1").
					WithArgs("notfound").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			url, err := pstStorage.GetURL(tt.shortURL)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedURL, url)
		})
	}

	// Проверяем, что все ожидания для mock были вызваны
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPstStorage_GetAllURL(t *testing.T) {

	// Создаем mock SQL соединение
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Инициализируем PstStorage с mock-базой данных
	pstStorage := &PstStorage{storage: db}

	tests := []struct {
		name         string
		userID       string
		baseURL      string
		expectedErr  error
		mockBehavior func()
	}{
		{
			name:        "successful",
			userID:      "test",
			baseURL:     "http://localhost:8080",
			expectedErr: nil,
			mockBehavior: func() {
				mock.ExpectBegin()
				rows := mock.NewRows([]string{"short_url", "original_url"}).
					AddRow("qwerty", "https://ya.ru")
				mock.ExpectQuery("SELECT short_url, original_url FROM urls WHERE user_id = \\$1").
					WithArgs("test").
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		{
			name:        "rows_error",
			userID:      "test",
			baseURL:     "http://localhost:8080",
			expectedErr: sql.ErrNoRows,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT short_url, original_url FROM urls WHERE user_id = \\$1").
					WithArgs("").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			_, err := pstStorage.GetAllURL(tt.userID, tt.baseURL)
			//if !errors.Is(err, tt.expectedErr) {
			//	t.Errorf("Ожидали ошибку = %v, а пришло = %v", tt.expectedErr, err)
			//}
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}

}
