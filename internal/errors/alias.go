package errors

import "errors"

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")
