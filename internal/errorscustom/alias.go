package errorscustom

import "errors"

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

// ErrUserIDNotContext указывает на отсутствие контекста userID.
var ErrUserIDNotContext = errors.New("userID not found or empty")

// ErrDeletedURL указывает на удаление URL.
var ErrDeletedURL = errors.New("URL DELETED")

// ErrBadVarifyToken указывает что токен не прошел верификацию
var ErrBadVarifyToken = errors.New("incorrect token")
