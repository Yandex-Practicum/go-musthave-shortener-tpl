package model


type URL struct {
	ID      string
	FullURL string
}

//Фабричный метод для создания экземпляра URL структуры
func NewURL(id, full string) *URL {
	return &URL{
		ID:      id,
		FullURL: full,
	}
}


type ApiPostRequest struct {
	URL string `json:"url"`
}

type ApiPostResponse struct {
	Result string `json:"result"`
}

func NewApiPostResponse (result string) *ApiPostResponse {
	return &ApiPostResponse{Result: result}
}