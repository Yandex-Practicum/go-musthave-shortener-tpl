package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
)

func (p *PstStorage) GetURL(shortURL string) (string, error) {
	var originalURL string
	db := p.storage
	// создаем запрос
	query := "SELECT original_url FROM urls WHERE short_url = $1"
	// делаем запрос
	row := db.QueryRowContext(context.Background(), query, shortURL)

	if row == nil {
		return "", sql.ErrNoRows
	}

	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}
	return originalURL, nil
}

func (p *PstStorage) GetAllURL(userID, baseURL string) ([]*models.UserURLs, error) {
	var userURLs []*models.UserURLs
	tx, err := p.storage.Begin()

	if err != nil {
		return nil, err
	}
	// создаем запрос
	query := "SELECT short_url, original_url FROM urls WHERE user_id = $1"

	// делаем запрос
	rows, err := tx.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, sql.ErrNoRows
	}
	defer rows.Close()

	//собираем все сохраненные ссылки от пользователя
	for rows.Next() {
		var userURL models.UserURLs
		if err = rows.Scan(&userURL.ShortURL, &userURL.OriginalURL); err != nil {
			return nil, err
		}
		userURL.ShortURL = fmt.Sprintf("%s/%s", baseURL, userURL.ShortURL)
		userURLs = append(userURLs, &userURL)

	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	// завершаем транзакцию
	tx.Commit()
	return userURLs, nil
}
