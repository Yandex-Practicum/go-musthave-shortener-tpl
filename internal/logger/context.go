package logger

import (
	"context"
	"log/slog"
)

type ctxLogger struct{}

// ContextWithLogger возвращает контекст с логгером.
func ContextWithLogger(ctx context.Context, sLogger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, sLogger)
}

// loggerFromContext возвращает логгер из контекста.
func loggerFromContext(ctx context.Context) *slog.Logger {
	if sLogger, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return sLogger
	}
	return slog.Default()
}
