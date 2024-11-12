package db

import (
	"context"

	"github.com/kamencov/go-musthave-shortener-tpl/internal/models"
	"github.com/kamencov/go-musthave-shortener-tpl/internal/utils"
)

// SaveURL сохраняет URL в базе данных.
func (p *PstStorage) SaveURL(shortURL, originalURL, userID string) error {
	var user *string
	if userID != "" {
		user = &userID
	}

	// начинаем транзакцию
	tx, err := p.storage.Begin()
	if err != nil {
		return err
	}
	// создаем запрос
	query := "INSERT INTO urls (original_url, short_url, user_id) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(context.Background(), query, originalURL, shortURL, user)
	if err != nil {
		// если ошибка, то откатываем изменения
		tx.Rollback()
		return err
	}

	// завершаем транзакцию
	return tx.Commit()
}

// SaveSliceOfDB сохраняет множество URL в базе данных.
func (p *PstStorage) SaveSlice(urls []models.MultipleURL, baseURL, userID string) ([]models.ResultMultipleURL, error) {
	var resultMultipleURL []models.ResultMultipleURL

	tx, err := p.storage.Begin()
	if err != nil {
		return resultMultipleURL, err
	}

	// создаем короткую ссылку и записываем в resultMultipleURL
	for _, req := range urls {
		// проверяем есть ли в базе уже данный URL
		encodeURL, err := p.CheckURL(req.OriginalURL)
		if err == nil {
			encodeURL, err = utils.EncodeURL(req.OriginalURL)
			if err != nil {
				return resultMultipleURL, err
			}
		}

		resultMultipleURL = append(resultMultipleURL, models.ResultMultipleURL{
			CorrelationID: req.CorrelationID,
			ShortURL:      baseURL + "/" + encodeURL,
		})

		p.SaveURL(encodeURL, req.OriginalURL, userID)
		tx.Rollback()
	}

	// завершаем транзакцию
	tx.Commit()
	return resultMultipleURL, nil
}
