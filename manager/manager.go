package manager

import (
	"context"
	"time"
)

var _ TokenManager = new(manager)

// KeyToken is the key to be used for storing the token.
const KeyToken = "manager:token:valuefirst"

// TokenManager is the contract of the token manager.
type TokenManager interface {
	// Get gets the token from the storage.Hub.
	Get(ctx context.Context) (token string, err error)
}

// Manager is a token manager.
type manager struct {
	opt *Option
}

// Assign assigns the option to the manager.
func (m *manager) Assign(o *Option) *manager {
	if o == nil {
		return m
	}

	m.opt = o
	return m
}

// New creates a new token manager.
func New(opts ...FnOption) (tm TokenManager) {
	return (new(manager)).Assign(
		(new(Option)).
			Assign(opts...).
			Default().
			Clone(),
	)
}

// Get gets the token from the storage.
func (m *manager) Get(ctx context.Context) (token string, err error) {
	err = m.opt.Validate()
	if err != nil {
		return
	}

	now := time.Now()
	res, err := m.opt.Storage.Get(ctx, m.opt.Key)
	if !(err != nil || res == nil || res.IsExpired(now)) {
		token = res.Data
		return
	}

	tokenResp, err := m.opt.Client.GenerateToken(ctx)
	if err != nil {
		return
	}

	_, err = m.opt.Client.EnableToken(ctx, tokenResp.Token)
	if err != nil {
		return
	}

	res, err = tokenResp.ToToken()
	if err != nil {
		return
	}

	err = m.opt.Storage.Save(ctx, m.opt.Key, res.SetHalfExpiredDate(now))
	if err != nil {
		return
	}

	token = res.Data
	return
}
