package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeURL(t *testing.T) {
	// Тестирование успешного случая
	t.Run("successful encoding", func(t *testing.T) {
		url := "https://example.com"
		expectedLen := 5
		encoded, err := EncodeURL(url)
		assert.Nil(t, err)
		assert.Equal(t, expectedLen, len(encoded), "Длина закодированного URL должна быть 6")
	})

	// Тестирование случая с пустым URL
	t.Run("empty URL", func(t *testing.T) {
		url := ""
		_, err := EncodeURL(url)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("URL is empty"), err, "Ошибка должна быть 'URL is empty'")
	})
}
