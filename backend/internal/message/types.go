package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type EventType int

const (
	EventTypeGroupChat EventType = iota + 1
	EventTypeGroupJoin
	EventTypeGroupLeave
	EventTypeGroupPromotion
)

type Message struct {
	GroupID  string    `json:"groupId"`
	PlayerID int       `json:"playerId"`
	Type     EventType `json:"type"`
	Payload  any       `json:"payload"`
}

func NewMessage(groupID string, playerID int, eventType EventType, payload interface{}) *Message {
	return &Message{
		GroupID:  groupID,
		PlayerID: playerID,
		Type:     eventType,
		Payload:  payload,
	}
}

func (e *Message) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func UnmarshalJSON(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %v", err)
	}
	return &msg, nil
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type Exchange interface {
	Publish(ctx context.Context, msg *Message) error
	Subscribe(ctx context.Context) *redis.PubSub // TODO: this should be more generic
}
