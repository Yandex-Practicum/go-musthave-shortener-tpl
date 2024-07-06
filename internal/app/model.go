package model


type URL struct {
	ID      string
	FullURL string
}

func NewURL(id, full string) *URL {
	return &URL{
		ID:      id,
		FullURL: full,
	}
}