package service

func (s *Service) Ping() error {
	return s.storage.Ping()
}
