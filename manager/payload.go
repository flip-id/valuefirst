package manager

import (
	"github.com/flip-id/valuefirst/storage"
	"time"
)

// ResponseGenerateToken is the response of GenerateToken.
type ResponseGenerateToken struct {
	Token string `json:"token"`

	// ExpiryDate is the expiry date of the token.
	// The format is "2022-07-10 22:32:30".
	// The time format used in this package is storage.TimeFormatExpiredDate.
	ExpiryDate string `json:"expiryDate"`
}

// ResponseEnableToken is the response of EnableToken.
type ResponseEnableToken struct {
	Response string `json:"Response"`
}

// ToToken converts the ResponseGenerateToken to a Token.
func (r *ResponseGenerateToken) ToToken() (s *storage.Token, err error) {
	newTime, err := time.Parse(storage.TimeFormatExpiredDate, r.ExpiryDate)
	if err != nil {
		return
	}

	s = &storage.Token{
		Data:        r.Token,
		ExpiredDate: newTime,
	}
	return
}
