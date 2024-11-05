package service

// Ping проверяет соединение с базой данных.
func (s *Service) Ping() error {
	return s.storage.Ping()
}
