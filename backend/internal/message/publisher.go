package message

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type IPublisher interface {
	PlayerJoined(ctx context.Context, groupID string, player *repository.PlayerInGroup) error
	PlayerLeft(ctx context.Context, groupID string, userID int, playerRemoved int, leaderID int) error
	GroupDeleted(ctx context.Context, groupID string, userID int) error
}

type Publisher struct {
	exchange Exchange
}

func NewPublisher(exchange Exchange) *Publisher {
	return &Publisher{
		exchange: exchange,
	}
}

func (p *Publisher) PlayerJoined(ctx context.Context, groupID string, player *repository.PlayerInGroup) error {
	msg := NewMessage(groupID, player.ID, EventTypeGroupJoin, player)
	return p.exchange.Publish(ctx, msg)
}

type PlayerLeftPayload struct {
	PlayerID int `json:"playerId"`
	LeaderID int `json:"leaderId"`
}

func (p *Publisher) PlayerLeft(ctx context.Context, groupID string, userID int, playerLeft int, leaderID int) error {
	msg := NewMessage(groupID, userID, EventTypeGroupLeave, &PlayerLeftPayload{
		PlayerID: playerLeft,
		LeaderID: leaderID,
	})
	return p.exchange.Publish(ctx, msg)
}

func (p *Publisher) GroupDeleted(ctx context.Context, groupID string, userID int) error {
	msg := NewMessage(groupID, userID, EventTypeGroupDelete, nil)
	return p.exchange.Publish(ctx, msg)
}
