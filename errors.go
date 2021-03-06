package valuefirst

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

// List of errors in the valuefirst package.
var (
	ErrEmptyUsername     = errors.New("empty username")
	ErrEmptyPassword     = errors.New("empty password")
	ErrDecodeVarIsNotPtr = errors.New("the decode variable is not a pointer")
	ErrNilRequest        = errors.New("nil request")
)

func formatUnknown(resp *http.Response) (err error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = &UnknownError{
		Message: string(b),
	}
	return
}

// UnknownError is an error that is not defined by the documentation from the ValueFirst.
type UnknownError struct {
	Message interface{} `json:"message"`
}

// Error implements the error interface.
func (e UnknownError) Error() string {
	return fmt.Sprintf(
		"unknown error ValueFirst: message: %v",
		e.Message,
	)
}

// ErrorList is a list of errors.
type ErrorList []error

// Error implements the error interface.
func (e ErrorList) Error() (res string) {
	byteSlice, _ := json.Marshal(e)
	res = string(byteSlice)
	return
}

// Append appends error to the ErrorList.
func (e *ErrorList) Append(err error) {
	*e = append(*e, err)
}
