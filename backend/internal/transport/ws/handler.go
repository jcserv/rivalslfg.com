package ws

import (
	"context"
	"encoding/json"
)

type DefaultHandler struct {
	hub *Hub
}

func NewDefaultHandler(hub *Hub) *DefaultHandler {
	return &DefaultHandler{hub: hub}
}

func (h *DefaultHandler) Handle(ctx context.Context, client *Client, payload json.RawMessage) error {
	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}
	return h.hub.Broadcast(ctx, msg)
}
