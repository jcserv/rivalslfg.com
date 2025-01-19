package ws

import (
	"context"
	"encoding/json"
)

type WebSocketEventType int

const (
	OpGroupChat      WebSocketEventType = iota + 1
	OpGroupJoin      WebSocketEventType = iota + 2
	OpGroupLeave     WebSocketEventType = iota + 3
	OpGroupPromotion WebSocketEventType = iota + 4
)

type EventHandler interface {
	Handle(ctx context.Context, client *Client, payload json.RawMessage) error
}
