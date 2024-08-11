package db

import "context"

func (p *PstStorage) SaveURL(shortURL, originalURL string) error {
	db := p.storage
	// начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// создаем запрос
	query := "INSERT INTO urls (originalURL, shortURL) VALUES ($1, $2)"
	_, err = tx.ExecContext(context.Background(), query, originalURL, shortURL)
	if err != nil {
		// если ошибка, то откатываем изменения
		tx.Rollback()
		return err
	}

	// завершаем транзакцию
	return tx.Commit()
}
