package db

import "context"

func (p *PstStorage) GetURL(shortURL string) (string, error) {
	var originalURL string
	db := p.storage
	// создаем запрос
	query := "SELECT originalURL FROM urls WHERE shortURL = $1"
	// делаем запрос
	row := db.QueryRowContext(context.Background(), query, shortURL)

	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}
	return originalURL, nil
}
