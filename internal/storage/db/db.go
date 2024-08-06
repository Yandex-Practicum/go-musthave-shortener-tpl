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

// InitDB инициализирует подключение к базе данных
func (p *PstStorage) initDB(dataSourceName string) error {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return err
	}
	p.storage = db
	fmt.Println(dataSourceName)
	fmt.Println(db)
	return nil
}

func (p *PstStorage) Ping() error {
	return p.storage.Ping()
}

func (p *PstStorage) Close() error {
	return p.storage.Close()
}
