package valuefirst

import (
	"encoding/json"
	"fmt"
	rh "github.com/fairyhunter13/reflecthelper/v5"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
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
func (e *ErrorList) Error() (res string) {
	byteSlice, _ := json.Marshal(e)
	res = string(byteSlice)
	return
}

// Append appends error to the ErrorList.
func (e *ErrorList) Append(err error) {
	if err == nil {
		return
	}

	*e = append(*e, err)
}

// ErrorCode is the error code from the ValueFirst.
type ErrorCode string

// Integer represents the integer type of ErrorCode.
func (e ErrorCode) Integer() (int, error) {
	res, err := strconv.Atoi(string(e))
	return res, err
}

// String returns the string representation of the ErrorCode.
func (e ErrorCode) String() string {
	return string(e)
}

// CastErrCode casts the input of interface{} to the ErrorCode type.
func CastErrCode(in interface{}) ErrorCode {
	return ErrorCode(rh.GetString(in))
}

// List of error codes in the ValueFirst.
const (
	// General

	ErrCodeGeneralSuccess       ErrorCode = "0"
	ErrCodeGeneralNotConformDTD ErrorCode = "65535"

	// Message Post

	ErrCodeMessageInvalidSenderID ErrorCode = "28680"
	ErrCodeMessageInvalidMessage  ErrorCode = "28681"

	// Status Request

	ErrCodeStatusRequestSuccess ErrorCode = "8448"

	// Scheduler Related

	ErrCodeSchedulerSuccess ErrorCode = "13568"
)

var ignoredErrorCode = map[ErrorCode]bool{
	ErrCodeGeneralSuccess:       true,
	ErrCodeStatusRequestSuccess: true,
	ErrCodeSchedulerSuccess:     true,
}

func filterError(in *ResponseMessageAckGUIDError) (err error) {
	if in == nil {
		return
	}

	_, ok := ignoredErrorCode[in.GetCode()]
	if ok {
		return
	}

	err = in
	return
}

func filterErrors(ins *ResponseMessageAckGUIDErrors) (err error) {
	if ins == nil {
		return
	}

	var (
		newList ResponseMessageAckGUIDErrors
		tempErr error
	)
	for _, in := range *ins {
		tempErr = filterError(&in)
		if tempErr == nil {
			continue
		}

		newList = append(newList, in)
	}

	if len(newList) <= 0 {
		return
	}

	err = &newList
	return
}
