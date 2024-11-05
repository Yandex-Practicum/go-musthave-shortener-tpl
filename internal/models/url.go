package models

// URL - структура для хранения URL.
type URL struct {
	URL string `json:"url"`
}

// ResultURL - структура для возвращения URL.
type ResultURL struct {
	URL string `json:"result"`
}

// MultipleURL - структура для хранения URL.
type MultipleURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ResultMultipleURL - структура для возвращения URL.
type ResultMultipleURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
