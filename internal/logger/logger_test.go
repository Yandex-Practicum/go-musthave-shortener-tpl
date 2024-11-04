package logger

import (
	"context"
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

// Тест для функции WithAddSource
func TestWithAddSource(t *testing.T) {
	// Начальные значения LoggerOptions
	options := &LoggerOptions{
		AddSource: false, // По умолчанию значение false
	}

	// Вызываем функцию WithAddSource с true
	optionFunc := WithAddSource(true)
	optionFunc(options)

	// Проверяем, что AddSource обновился до true
	if !options.AddSource {
		t.Errorf("expected AddSource to be true, got %v", options.AddSource)
	}

	// Вызываем функцию WithAddSource с false
	optionFunc = WithAddSource(false)
	optionFunc(options)

	// Проверяем, что AddSource обновился до false
	if options.AddSource {
		t.Errorf("expected AddSource to be false, got %v", options.AddSource)
	}
}

// Тест для функции WithIsJSON
func TestWithIsJSON(t *testing.T) {
	// Начальные значения LoggerOptions
	options := &LoggerOptions{
		IsJSON: false, // По умолчанию значение false
	}

	// Вызываем функцию WithIsJSON с true
	optionFunc := WithIsJSON(true)
	optionFunc(options)

	// Проверяем, что IsJSON обновился до true
	if !options.IsJSON {
		t.Errorf("expected IsJSON to be true, got %v", options.IsJSON)
	}

	// Вызываем функцию WithIsJSON с false
	optionFunc = WithIsJSON(false)
	optionFunc(options)

	// Проверяем, что IsJSON обновился до false
	if options.IsJSON {
		t.Errorf("expected IsJSON to be false, got %v", options.IsJSON)
	}
}

// Тест для функции WithSetDefault
func TestWithSetDefault(t *testing.T) {
	// Начальные значения LoggerOptions
	options := &LoggerOptions{
		SetDefault: false, // По умолчанию значение false
	}

	// Вызываем функцию WithSetDefault с true
	optionFunc := WithSetDefault(true)
	optionFunc(options)

	// Проверяем, что SetDefault обновился до true
	if !options.SetDefault {
		t.Errorf("expected SetDefault to be true, got %v", options.SetDefault)
	}

	// Вызываем функцию WithSetDefault с false
	optionFunc = WithSetDefault(false)
	optionFunc(options)

	// Проверяем, что SetDefault обновился до false
	if options.SetDefault {
		t.Errorf("expected SetDefault to be false, got %v", options.SetDefault)
	}
}

// Тест для функции WithAttr
func TestWithAttr(t *testing.T) {
	ctx := context.Background()
	logger := WithAttr(ctx, Attr{
		Key:   "key",
		Value: Value{},
	})
	if logger == nil {
		t.Errorf("expected logger to be not nil, got nil")
	}
}

// Тест для функции WithDefaultAttrs
func TestWithDefaultAttrs(t *testing.T) {
	logger := WithDefaultAttrs(Default(), Attr{
		Key:   "key",
		Value: Value{},
	})
	if logger == nil {
		t.Errorf("expected logger to be not nil, got nil")
	}
}
