package filestorage

import "errors"

func (s *SaveFile) DeletedURLs(doneCh chan struct{}, urlCh chan string, userID string) error {
	return errors.New("there is no such implementation")
}
