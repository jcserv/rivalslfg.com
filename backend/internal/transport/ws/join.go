package ws

import (
	"context"
	"encoding/json"
)

type JoinPayload struct {
	Name string `json:"name"`
}

type JoinHandler struct {
	hub *Hub
}

func NewJoinHandler(hub *Hub) *JoinHandler {
	return &JoinHandler{hub: hub}
}

func (h *JoinHandler) Handle(ctx context.Context, client *Client, payload json.RawMessage) error {
	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}
	return h.hub.Broadcast(ctx, msg)
}
