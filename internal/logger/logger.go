package logger

import (
	"context"
	"log/slog"
	"os"
)

//
//func Initialize(level string) error {
//	slogLevel, err := purseLevel(level)
//	if err != nil {
//		return err
//	}
//
//	// указываем параметры
//	opts := &slog.HandlerOptions{
//		Level: slogLevel,
//	}
//
//	// указываем тип хендлера - вывод + параметры
//	handler := slog.NewJSONHandler(os.Stdout, opts)
//
//	// создаем новый логер
//	logger := slog.New(handler)
//
//	// делаем его логером по умолчанию
//	slog.SetDefault(logger)
//	return nil
//}
//
//func purseLevel(level string) (slog.Level, error) {
//	level = strings.ToLower(level)
//	switch level {
//	case "debug":
//		return slog.LevelDebug, nil
//	case "info":
//		return slog.LevelInfo, nil
//	case "warn":
//		return slog.LevelWarn, nil
//	case "error":
//		return slog.LevelError, nil
//	default:
//		return slog.Level(999), errors.New("Bad flag Log Level")
//	}
//}

const (
	defaultLevel      = LevelInfo
	defaultAddSource  = false
	defaultIsJSON     = true
	defaultSetDefault = true
)

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

func WithLevel(level string) LoggerOprion {
	return func(optns *LoggerOptions) {
		var lvl Level
		if err := lvl.UnmarshalText([]byte(level)); err != nil {
			lvl = LevelInfo
		}

		optns.Level = lvl
	}
}

func WithAddSource(addSource bool) LoggerOprion {
	return func(optns *LoggerOptions) {
		optns.AddSource = addSource
	}
}

func WithIsJSON(isJSON bool) LoggerOprion {
	return func(optns *LoggerOptions) {
		optns.IsJSON = isJSON
	}
}

func WithSetDefault(setDefault bool) LoggerOprion {
	return func(optns *LoggerOptions) {
		optns.SetDefault = setDefault
	}
}

func WithAttr(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

func WithDefaultAttrs(logger *Logger, attrs ...Attr) *Logger {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}

func Default() *Logger {
	return slog.Default()
}
