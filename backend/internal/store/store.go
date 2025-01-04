package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any) error
	// Delete(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, duration time.Duration) error
}

func New(conn *redis.Client) Store {
	return NewRedisStore(conn)
}
