package service

type Storage interface {
	SaveURL(string, string) error
	GetURL(string) (string, error)
	Close() error
	Ping() error
	CheckURL(string) (string, error)
}
