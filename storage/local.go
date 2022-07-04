package storage

import (
	"context"
	"fmt"
	"sync"
)

var _ Hub = new(localStorage)

type localStorage struct {
	sync.Map
}

// NewLocalStorage creates a new local storage to store the token.
func NewLocalStorage() Hub {
	return &localStorage{}
}

// Get gets the token from the memory storage.
func (s *localStorage) Get(ctx context.Context, key string) (t *Token, err error) {
	v, ok := s.Load(key)
	if !ok {
		err = fmt.Errorf("token is not found, key: %s", key)
		return
	}

	t, ok = v.(*Token)
	if !ok {
		err = fmt.Errorf("token is not a valid Token struct, key: %s", key)
	}
	return
}

// Save saves the token to the memory storage.
func (s *localStorage) Save(ctx context.Context, key string, t *Token) (err error) {
	if t == nil {
		err = ErrNilToken
		return
	}

	s.Store(key, t)
	return
}
