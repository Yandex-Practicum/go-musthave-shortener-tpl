package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	defaultLevel      = LevelInfo
	defaultAddSource  = false
	defaultIsJSON     = true
	defaultSetDefault = true
)

// NewLogger создает новый логгер.
func NewLogger(opts ...LoggerOprion) *Logger {
	config := &LoggerOptions{
		Level:      defaultLevel,
		AddSource:  defaultAddSource,
		IsJSON:     defaultIsJSON,
		SetDefault: defaultSetDefault,
	}

	for _, opt := range opts {
		opt(config)
	}

	options := &HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	var h Handler = NewTextHandler(os.Stdout, options)

	if config.IsJSON {
		h = NewJSONHandler(os.Stdout, options)
	}

	logger := New(h)

	if config.SetDefault {
		SetDefault(logger)
	}

	return logger
}

type LoggerOptions struct {
	Level      Level
	AddSource  bool
	IsJSON     bool
	SetDefault bool
}

type LoggerOprion func(options *LoggerOptions)

// WithLevel устанавливает уровень логирования.
func WithLevel(level string) LoggerOprion {
	return func(optns *LoggerOptions) {
		var lvl Level
		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			lvl = LevelInfo
		}

		optns.Level = lvl
	}
}

// WithAddSource добавляет источник логирования.
func WithAddSource(addSource bool) LoggerOprion {
	return func(optns *LoggerOptions) {
		optns.AddSource = addSource
	}
}

// WithIsJSON добавляет источник логирования.
func WithIsJSON(isJSON bool) LoggerOprion {
	return func(optns *LoggerOptions) {
		optns.IsJSON = isJSON
	}
}

// WithSetDefault добавляет источник логирования.
func WithSetDefault(setDefault bool) LoggerOprion {
	return func(optns *LoggerOptions) {
		optns.SetDefault = setDefault
	}
}

// WithAttr устанавливает аттрибуты логирования.
func WithAttr(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

// WithDefaultAttrs устанавливает аттрибуты логирования.
func WithDefaultAttrs(logger *Logger, attrs ...Attr) *Logger {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

// L возвращает логгер из контекста.
func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}

// Default возвращает логгер по умолчанию.
func Default() *Logger {
	return slog.Default()
}
