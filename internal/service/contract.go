package service

type Storage interface {
	SaveURL(string, string) error
	GetURL(string) (string, error)
}
