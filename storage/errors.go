package storage

import "github.com/pkg/errors"

var (
	// ErrNilToken specifies that the token passed to the function is nil.
	ErrNilToken = errors.New("token is nil")
)
