package mapstorage

import "errors"

func (s *MapStorage) DeletedURLs(doneCh chan struct{}, urlCh chan string, userID string) error {
	return errors.New("there is no such implementation")
}
