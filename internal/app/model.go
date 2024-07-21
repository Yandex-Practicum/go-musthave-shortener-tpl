package model

type URL struct {
	ID      string
	FullURL string
}

// Фабричный метод для создания экземпляра URL структуры
func NewURL(id, full string) *URL {
	return &URL{
		ID:      id,
		FullURL: full,
	}
}

type APIPostRequest struct {
	URL string `json:"url"`
}

type APIPostResponse struct {
	Result string `json:"result"`
}

func NewAPIPostResponse(result string) *APIPostResponse {
	return &APIPostResponse{Result: result}
}
