package logger

import (
	"context"
	"testing"
)

// Тест для функции ContextWithLogger
func TestContextWithLogger(t *testing.T) {
	logger := NewLogger()
	ctx := ContextWithLogger(context.Background(), logger)

	if ctx == nil {
		t.Errorf("expected context to be not nil, got nil")
	}
}
