package url

import "sync"

type mapUrl struct {
	Storage map[string]string
	mu      *sync.RWMutex
}

func NewMapUrl() *mapUrl {
	return &mapUrl{
		Storage: make(map[string]string),
	}
}
