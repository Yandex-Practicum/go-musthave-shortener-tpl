package models

type URL struct {
	URL string `json:"url"`
}

type ResultURL struct {
	URL string `json:"result"`
}

type MultipleURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResultMultipleURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
