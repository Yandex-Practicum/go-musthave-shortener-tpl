package service

func (s *Service) DeletedURLs(url []string, userID string) error {
	return s.storage.DeletedURLs(url, userID)
}
