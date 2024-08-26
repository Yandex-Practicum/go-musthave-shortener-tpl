package mapstorage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapStorage_SaveURL(t *testing.T) {
	t.Run("successful_saving", func(t *testing.T) {
		s := NewMapURL()
		err := s.SaveURL("test", "", "")
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("URL is empty"), err)
		err = s.SaveURL("test", "https://example.com", "")
		assert.Nil(t, err)
	})
}

func TestMapStorage_GetURL(t *testing.T) {
	t.Run("successful_getting", func(t *testing.T) {
		s := NewMapURL()
		err := s.SaveURL("test", "https://example.com", "")
		assert.Nil(t, err)
		_, err = s.GetURL("")
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("URL not found"), err)
		url, err := s.GetURL("test")
		assert.Nil(t, err)
		assert.Equal(t, "https://example.com", url)
	})
}
