package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestNewPstStorage(t *testing.T) {
	// Создаем mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания sqlmock: %v", err)
	}
	defer db.Close()

	// Определяем поведение mock для операций с базой данных
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS urls").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("ALTER TABLE urls ADD COLUMN IF NOT EXISTS is_deleted").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Создаем тестовое хранилище
	storage := &PstStorage{storage: db}

	// Вызываем метод CreateTableIfNotExists
	err = storage.CreateTableIfNotExists()
	if err != nil {
		t.Errorf("ошибка создания таблицы: %v", err)
	}

	// Проверяем, что все ожидаемые действия были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не выполнены все ожидания mock: %v", err)
	}
}
