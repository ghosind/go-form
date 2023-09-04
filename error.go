package form

import "errors"

var (
	// ErrUnsupportedType indicates the type of value that to be encoded is unsupported.
	ErrUnsupportedType error = errors.New("unsupported value type")
)
