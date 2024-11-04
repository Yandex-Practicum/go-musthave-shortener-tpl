package logger

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

// Тест для функции Float32Attr
func TestFloat32Attr(t *testing.T) {
	attr := Float32Attr("float32_test", float32(1.23))
	expected := slog.Float64("float32_test", 1.23)

	if attr.Key != expected.Key {
		t.Errorf("expected %v, got %v", expected, attr)
	}
}

// Тест для функции Uint32Attr
func TestUint32Attr(t *testing.T) {
	attr := Uint32Attr("uint32_test", uint32(123))
	expected := slog.Int("uint32_test", 123)

	if attr.Key != expected.Key {
		t.Errorf("expected %v, got %v", expected, attr)
	}
}

// Тест для функции Int32Attr
func TestInt32Attr(t *testing.T) {
	attr := Int32Attr("int32_test", int32(123))
	expected := slog.Int("int32_test", 123)

	if attr.Key != expected.Key {
		t.Errorf("expected %v, got %v", expected, attr)
	}
}

// Тест для функции TimeAttr
func TestTimeAttr(t *testing.T) {
	now := time.Now()
	attr := TimeAttr("time_test", now)
	expected := slog.String("time_test", now.String())

	if attr.Key != expected.Key {
		t.Errorf("expected %v, got %v", expected, attr)
	}
}

// Тест для функции ErrAttr
func TestErrAttr(t *testing.T) {
	err := errors.New("test error")
	attr := ErrAttr(err)
	expected := slog.String("error", err.Error())

	if attr.Key != expected.Key {
		t.Errorf("expected %v, got %v", expected, attr)
	}
}
