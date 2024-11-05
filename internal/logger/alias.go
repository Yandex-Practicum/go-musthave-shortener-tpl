package logger

import (
	"log/slog"
	"time"
)

// уровни логирования.
const (
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelDebug = slog.LevelDebug
)

// псевдонимы для slog.
type (
	Logger         = slog.Logger
	Attr           = slog.Attr
	Level          = slog.Level
	Handler        = slog.Handler
	Value          = slog.Value
	HandlerOptions = slog.HandlerOptions
	LogValuer      = slog.LogValuer
)

// псевдонимы для slog.
var (
	NewTextHandler = slog.NewTextHandler
	NewJSONHandler = slog.NewJSONHandler
	New            = slog.New
	SetDefault     = slog.SetDefault

	StringAttr   = slog.String
	BoolAttr     = slog.Bool
	Float64Attr  = slog.Float64
	AnyAttr      = slog.Any
	DurationAttr = slog.Duration
	IntAttr      = slog.Int
	Int64Attr    = slog.Int64
	Uint64Attr   = slog.Uint64

	GroupValue = slog.GroupValue
	Group      = slog.Group
)

// Float32Attr преобразует float32 в float64.
func Float32Attr(key string, val float32) Attr {
	return slog.Float64(key, float64(val))
}

// Uint32Attr преобразует uint32 в int.
func Uint32Attr(key string, val uint32) Attr {
	return slog.Int(key, int(val))
}

// Int32Attr преобразует int32 в int.
func Int32Attr(key string, val int32) Attr {
	return slog.Int(key, int(val))
}

// TimeAttr преобразует time.Time в string.
func TimeAttr(key string, time time.Time) Attr {
	return slog.String(key, time.String())
}

// ErrAttr преобразует error в string.
func ErrAttr(err error) Attr {
	return slog.String("error", err.Error())
}
