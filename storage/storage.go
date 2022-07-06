package storage

import (
	"context"
	"time"
)

// List of constants used in this package.
const (
	// TimeFormatExpiredDate is the default time layout for golang time.Parse.
	TimeFormatExpiredDate = "2006-01-02 15:04:05"
)

// Token specifies the token data.
type Token struct {
	Data        string
	Duration    time.Duration
	ExpiredDate time.Time
}

// SetHalfExpiredDate sets the expired date to the half of the duration.
func (t *Token) SetHalfExpiredDate(now time.Time) *Token {
	if t.IsExpired(now) {
		return t
	}

	halfDuration := t.ExpiredDate.Sub(now) / 2
	return t.reduceExpiredDate(halfDuration).adjustDuration(now)
}

func (t *Token) adjustDuration(now time.Time) *Token {
	t.Duration = t.ExpiredDate.Sub(now)
	return t
}

func (t *Token) reduceExpiredDate(d time.Duration) *Token {
	if d < 0 {
		t.ExpiredDate = t.ExpiredDate.Add(d)
	} else {
		t.ExpiredDate = t.ExpiredDate.Add(-1 * d)
	}
	return t
}

// IsExpired checks whether the token is expired.
func (t *Token) IsExpired(now time.Time) bool {
	return now == time.Time{} || t.ExpiredDate.Before(now) || t.ExpiredDate.Equal(now)
}

// Hub specifies the contract to interact with the storage provider.
type Hub interface {
	Get(ctx context.Context, key string) (token *Token, err error)
	Save(ctx context.Context, key string, token *Token) (err error)
}
