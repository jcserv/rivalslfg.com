package store

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	conn *redis.Client
}

func NewRedisStore(conn *redis.Client) *RedisStore {
	return &RedisStore{
		conn: conn,
	}
}

func (s *RedisStore) Get(ctx context.Context, key string) (any, error) {
	result, err := s.conn.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, errors.New("key not found")
	}
	return result, nil
}

func (s *RedisStore) Set(ctx context.Context, key string, value any) error {
	return s.conn.HSet(ctx, key, value).Err()
}

func (s *RedisStore) Expire(ctx context.Context, key string, duration time.Duration) error {
	return s.conn.Expire(ctx, key, duration).Err()
}
