package db

import (
	"fmt"
)

// DeletedURLs удаляет URL из базы данных.
func (p *PstStorage) DeletedURLs(urls []string, userID string) error {
	if len(urls) == 0 {
		return nil
	}

	// Создаем SQL-запрос для множественного обновления.
	query := `UPDATE urls SET is_deleted = TRUE WHERE user_id = $1 AND short_url = ANY($2)`

	// Подготавливаем список URL в формате PostgreSQL.
	urlsArray := "{"
	for i, url := range urls {
		if i > 0 {
			urlsArray += ","
		}
		urlsArray += fmt.Sprintf(`"%s"`, url)
	}
	urlsArray += "}"

	// Выполняем запрос.
	_, err := p.storage.Exec(query, userID, urlsArray)
	if err != nil {
		return err
	}

	return nil
}
