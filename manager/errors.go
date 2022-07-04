package manager

import "github.com/pkg/errors"

var (
	// ErrNilClient is the error for nil client.
	ErrNilClient = errors.New("nil client")
)
