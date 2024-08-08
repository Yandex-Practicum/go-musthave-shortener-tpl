package db

import "context"

func (p *PstStorage) SaveURL(shortURL, originalURL string) error {
	db := p.storage

	// создаем запрос
	query := "INSERT INTO urls (originalURL, shortURL) VALUES ($1, $2)"
	_, err := db.ExecContext(context.Background(), query, originalURL, shortURL)
	if err != nil {
		return err
	}
	return nil
}
