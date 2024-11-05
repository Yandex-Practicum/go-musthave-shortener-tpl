package filestorage

import "errors"

// DeletedURLs удаляет URL из файла.
func (s *SaveFile) DeletedURLs(url []string, userID string) error {
	return errors.New("there is no such implementation")
}
