package service

//go:generate mockgen -source=./contract.go -destination=../mocks/mock_storage.go -package=mocks
type Storage interface {
	SaveURL(string, string) error
	GetURL(string) (string, error)
	Close() error
	Ping() error
	CheckURL(string) (string, error)
}
