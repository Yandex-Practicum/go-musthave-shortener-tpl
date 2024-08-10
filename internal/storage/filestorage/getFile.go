package filestorage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

func (s *SaveFile) GetURL(shortURL string) (string, error) {
	// Читаем содержимое файла
	readFile, err := os.Open(s.file.Name())
	if err != nil {
		return "", err
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	for scanner.Scan() {
		var event Event
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &event)
		if err != nil {
			continue // Пропуск некорректных JSON строк
		}
		if event.ShortURL == shortURL {
			return event.OriginalURL, nil
		}
	}

	return "", errors.New("короткий URL не найден")
}
