package service

// DeletedURLs - удаление URL из хранилища
func (s *Service) DeletedURLs(url []string, userID string) error {
	return s.storage.DeletedURLs(url, userID)
}
