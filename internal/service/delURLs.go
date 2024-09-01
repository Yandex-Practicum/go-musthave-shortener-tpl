package service

func (s *Service) DeletedURLs(doneCh chan struct{}, urlCh chan string, userID string) error {
	return s.storage.DeletedURLs(doneCh, urlCh, userID)
}
