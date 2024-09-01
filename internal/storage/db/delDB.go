package db

import (
	"database/sql"
	"fmt"
)

func (p *PstStorage) DeletedURLs(doneCh chan struct{}, urlCh chan string, userID string) error {
	var urls []string
	updated := false
	for {
		select {
		case url, ok := <-urlCh:
			if !ok {
				// Канал закрыт
				if !updated && len(urls) > 0 {
					err := updateDatabase(p.storage, userID, urls)
					if err != nil {
						return err
					}
				}
				return nil
			}
			urls = append(urls, url)
		case <-doneCh:
			// Получен сигнал завершения работы
			if len(urls) > 0 {
				err := updateDatabase(p.storage, userID, urls)
				if err != nil {
					return err
				}
				updated = true
			}
			return nil
		}
	}
}

// updateDatabase выполняет пакетное обновление в базе данных.
func updateDatabase(db *sql.DB, userID string, urls []string) error {
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
	_, err := db.Exec(query, userID, urlsArray)
	if err != nil {
		return err
	}

	return nil
}
