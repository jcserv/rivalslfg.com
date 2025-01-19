package message

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

const (
	GroupUpdateChannel = "group_updates"
)

type RedisExchange struct {
	conn *redis.Client
}

func NewRedisExchange(conn *redis.Client) *RedisExchange {
	return &RedisExchange{
		conn: conn,
	}
}

func (e *RedisExchange) Publish(ctx context.Context, msg *Message) error {
	err := e.conn.Publish(ctx, GroupUpdateChannel, msg).Err()
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error publishing message to Redis: %v", err))
	}
	return err
}

func (e *RedisExchange) Subscribe(ctx context.Context) (*redis.PubSub, error) {
	pubsub := e.conn.Subscribe(ctx, GroupUpdateChannel)
	_, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		return nil, err
	}
	return pubsub, nil
}
