package valuefirst

import (
	"github.com/fairyhunter13/reflecthelper/v5"
	"github.com/flip-id/valuefirst/manager"
	"github.com/flip-id/valuefirst/storage"
	"net/http"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
)

const (
	// DefaultTimeout sets the default timeout of the HTTP client.
	DefaultTimeout = 30 * time.Second

	// BaseURL sets the base URL of the API.
	BaseURL = "https://api.myvfirst.com/psms"
)

// FnOption is a function that sets the option.
type FnOption func(o *Option)

// WithBaseURL sets the base URL of the API.
func WithBaseURL(s string) FnOption {
	return func(o *Option) {
		o.BaseURL = s
	}
}

// WithBasicAuth sets the basic auth credentials.
func WithBasicAuth(user string, password string) FnOption {
	return func(o *Option) {
		o.BasicAuth = OptionBasicAuth{
			User:     user,
			Password: password,
		}
	}
}

// WithTimeout sets the timeout of the HTTP client.
func WithTimeout(t time.Duration) FnOption {
	return func(o *Option) {
		o.Timeout = t
	}
}

// WithClient sets the HTTP client.
func WithClient(c heimdall.Doer) FnOption {
	return func(o *Option) {
		o.Client = c
	}
}

// WithHystrixOptions sets the hystrix options.
func WithHystrixOptions(opts ...hystrix.Option) FnOption {
	return func(o *Option) {
		o.HystrixOptions = append(o.HystrixOptions, opts...)
	}
}

// WithCustomIPs sets the custom IPs used for hitting the API.
func WithCustomIPs(ips ...string) FnOption {
	return func(o *Option) {
		o.CustomIPs = ips
	}
}

// WithTokenStorage sets the token storage.
func WithTokenStorage(s storage.Hub) FnOption {
	return func(o *Option) {
		o.TokenStorage = s
	}
}

// WithTokenManager sets the token manager.
func WithTokenManager(tm manager.TokenManager) FnOption {
	return func(o *Option) {
		o.TokenManager = tm
	}
}

// Option is a config for Value.
type Option struct {
	BaseURL        string
	BasicAuth      OptionBasicAuth
	Timeout        time.Duration
	Client         heimdall.Doer
	HystrixOptions []hystrix.Option
	CustomIPs      []string
	TokenStorage   storage.Hub
	TokenManager   manager.TokenManager
	client         *hystrix.Client
	vfClient       *client
}

// OptionBasicAuth is a config for basic authorization.
type OptionBasicAuth struct {
	User     string
	Password string
}

// Clone clones the Option.
// Clone only makes a shallow copy of the Option struct.
func (o *Option) Clone() *Option {
	opt := *o
	return &opt
}

// Assign assigns the Option using the functional options.
func (o *Option) Assign(opts ...FnOption) *Option {
	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (o *Option) setValueFirstClient(c *client) *Option {
	o.vfClient = c
	return o
}

// Default sets the config default value.
func (o *Option) Default() *Option {
	if o.BaseURL == "" {
		o.BaseURL = BaseURL
	}

	o.BaseURL = strings.TrimRight(o.BaseURL, "/")
	if o.Timeout < DefaultTimeout {
		o.Timeout = DefaultTimeout
	}

	if o.Client == nil {
		o.Client = http.DefaultClient
	}

	if reflecthelper.IsNil(o.TokenStorage) {
		o.TokenStorage = storage.NewLocalStorage()
	}

	o.client = hystrix.NewClient(
		append(
			[]hystrix.Option{
				hystrix.WithHTTPTimeout(o.Timeout),
				hystrix.WithHystrixTimeout(o.Timeout),
				hystrix.WithHTTPClient(o.Client),
			},
			o.HystrixOptions...,
		)...,
	)

	if o.vfClient == nil {
		o.vfClient = &client{o}
	}

	if reflecthelper.IsNil(o.TokenManager) {
		o.TokenManager = manager.New(
			manager.WithClient(o.vfClient),
			manager.WithStorage(o.TokenStorage),
		)
	}
	return o
}

// Validate validates the config variables to ensure smooth integration.
func (o *Option) Validate() (err error) {
	if o.BasicAuth.User == "" {
		err = ErrEmptyUsername
		return
	}

	if o.BasicAuth.Password == "" {
		err = ErrEmptyPassword
	}
	return
}
