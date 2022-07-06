package storage

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

// Static Check
var _ Hub = new(redisStorage)

// NewGoRedisV8 creates a new redis client for storage.Hub.
func NewGoRedisV8(c *redis.Client) (h Hub) {
	if c == nil {
		return nil
	}

	h = &redisStorage{
		Client: c,
	}
	return
}

type redisStorage struct {
	*redis.Client
}

// Get returns the token from the Redis.
func (r *redisStorage) Get(ctx context.Context, key string) (t *Token, err error) {
	cmd := r.Client.Get(ctx, key)
	err = cmd.Err()
	if err != nil {
		return
	}

	byteSlice, err := cmd.Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(byteSlice, &t)
	return
}

// Save saves the token to the Redis storage.
func (r *redisStorage) Save(ctx context.Context, key string, t *Token) (err error) {
	if t == nil {
		err = ErrNilToken
		return
	}

	byteSlice, err := json.Marshal(t)
	if err != nil {
		return
	}

	cmd := r.Client.SetEX(ctx, key, string(byteSlice), t.Duration)
	err = cmd.Err()
	return
}
