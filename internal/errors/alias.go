package errors

import "errors"

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")
var ErrUserIDNotContext = errors.New("userID not found or empty")
var ErrDeletedURL = errors.New("URL DELETED")
