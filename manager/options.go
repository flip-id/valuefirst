package manager

import (
	"github.com/fairyhunter13/reflecthelper/v5"
	"github.com/flip-id/valuefirst/storage"
)

// Option is the option of the manager.
type Option struct {
	Storage storage.Hub
	Client  TokenClient
	Key     string
}

// Assign assigns the option.
func (o *Option) Assign(opts ...FnOption) *Option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Clone clones the Option.
func (o *Option) Clone() *Option {
	opt := *o
	return &opt
}

// FnOption is the function to set the option.
type FnOption func(o *Option)

// WithStorage sets the storage.
func WithStorage(s storage.Hub) FnOption {
	return func(o *Option) {
		o.Storage = s
	}
}

// WithClient sets the client.
func WithClient(c TokenClient) FnOption {
	return func(o *Option) {
		o.Client = c
	}
}

// WithKey sets the key.
func WithKey(key string) FnOption {
	return func(o *Option) {
		o.Key = key
	}
}

// Validate validates the option.
func (o *Option) Validate() (err error) {
	if reflecthelper.IsNil(o.Client) {
		err = ErrNilClient
	}
	return
}

// Default sets the default option.
func (o *Option) Default() *Option {
	if o.Key == "" {
		o.Key = KeyToken
	}

	if reflecthelper.IsNil(o.Storage) {
		o.Storage = storage.NewLocalStorage()
	}
	return o
}
