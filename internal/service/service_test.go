package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("create_service", func(t *testing.T) {
		s := NewService(nil, nil)
		assert.NotNil(t, s)
	})
}
