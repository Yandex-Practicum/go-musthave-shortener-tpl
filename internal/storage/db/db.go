package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PsqlStorage interface {
	initDB(dataSourceName string) error
	Ping() error
	Close() error
}

type PstStorage struct {
	storage *sql.DB
}

func NewPstStorage(dataSourceName string) (*PstStorage, error) {
	p := &PstStorage{}
	err := p.initDB(dataSourceName)
	return p, err
}

// InitDB инициализирует подключение к базе данных и создаем базу
func (p *PstStorage) initDB(dataSourceName string) error {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return err
	}
	p.storage = db
	fmt.Println(dataSourceName)
	fmt.Println(db)

	err = p.CreateTableIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

// Функция для создания таблицы, если она не существует
func (p *PstStorage) CreateTableIfNotExists() error {
	db, err := p.storage.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()

	queryCreate := `
    CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        original_url TEXT NOT NULL,
        short_url TEXT NOT NULL,
        user_id UUID,
        UNIQUE (original_url)
    );`

	_, err = db.Exec(queryCreate)
	if err != nil {
		return err
	}

	queryUpdate := `
    ALTER TABLE urls ADD COLUMN IF NOT EXISTS is_deleted BOOL NOT NULL DEFAULT FALSE;`

	_, err = db.Exec(queryUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (p *PstStorage) Ping() error {
	return p.storage.Ping()
}

func (p *PstStorage) Close() error {
	return p.storage.Close()
}
