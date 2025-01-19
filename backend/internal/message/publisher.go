package message

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type IPublisher interface {
	JoinGroup(ctx context.Context, groupID string, player *repository.PlayerInGroup) error
}

type Publisher struct {
	exchange Exchange
}

func NewPublisher(exchange Exchange) *Publisher {
	return &Publisher{
		exchange: exchange,
	}
}

func (p *Publisher) JoinGroup(ctx context.Context, groupID string, player *repository.PlayerInGroup) error {
	msg := NewMessage(groupID, player.ID, EventTypeGroupJoin, player)
	return p.exchange.Publish(ctx, msg)
}
