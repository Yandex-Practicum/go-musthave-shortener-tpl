package db

import (
	"context"
	"database/sql"
	"errors"

	errors2 "github.com/kamencov/go-musthave-shortener-tpl/internal/errorscustom"
)

func (p *PstStorage) CheckURL(originalURL string) (string, error) {
	var shortURL string

	if row := p.storage.QueryRowContext(context.Background(),
		"SELECT short_url FROM urls WHERE original_url = $1",
		originalURL).Scan(&shortURL); errors.Is(row, sql.ErrNoRows) {
		return "", nil
	}

	return shortURL, errors2.ErrConflict
}

func (p *PstStorage) CheckUser(user string) error {

	if row := p.storage.QueryRowContext(context.Background(),
		"SELECT user_id FROM urls WHERE user_id = $1",
		user); row.Err() != nil {
		return row.Err()
	}
	return nil
}
