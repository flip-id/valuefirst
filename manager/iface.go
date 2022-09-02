package manager

//go:generate mockgen -source=iface.go -destination=../mocks/iface.go -package=mocks

import "context"

// TokenClient is the contract of the token client to manage token.
type TokenClient interface {
	// GenerateToken generate a token for the user.
	GenerateToken(ctx context.Context) (res ResponseGenerateToken, err error)
	// EnableToken enable the token.
	EnableToken(ctx context.Context, token string) (res ResponseEnableToken, err error)
	// DisableToken disables the token.
	DisableToken(ctx context.Context, token string) (res ResponseEnableToken, err error)
	// DeleteToken disables the token.
	DeleteToken(ctx context.Context, token string) (res ResponseEnableToken, err error)
}
