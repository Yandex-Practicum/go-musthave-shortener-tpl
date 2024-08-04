package db

import "database/sql"

type DB struct {
	DB *sql.DB
}

// GetDB возвращает экземпляр базы данных
func NewDB() *DB {
	return &DB{}
}

// InitDB инициализирует подключение к базе данных
func InitDB(dataSourceName string) error {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return err
	}

	defer db.Close()
	return err
}
