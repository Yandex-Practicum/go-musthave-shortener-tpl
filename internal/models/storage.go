package models

// Storage - структура для хранения в базе данных.
type Storage struct {
	UUID        string `db:"user_id" json:"user_id"`
	ShortURL    string `db:"short_url" json:"short_url"`
	OriginalURL string `db:"original_url" json:"original_url"`
	DeletedFlag bool   `db:"is_deleted" json:"is_deleted"`
}
