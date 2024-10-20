package logger

import (
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *Logger
	}{
		{
			name: "test_default",
			want: Default(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithLevel(t *testing.T) {
	tests := []struct {
		input       string
		expectedLvl Level
	}{
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"debug", LevelDebug},
		{"error", LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			opts := &LoggerOptions{}
			WithLevel(tt.input)(opts)

			if !reflect.DeepEqual(opts.Level, tt.expectedLvl) {
				t.Errorf("Expected opts.Level to be %v, got %v", tt.expectedLvl, opts.Level)
			}
		})
	}
}
