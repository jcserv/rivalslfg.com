package ws

import (
	"context"
	"encoding/json"
)

type LeaveHandler struct {
	hub *Hub
}

func NewLeaveHandler(hub *Hub) *JoinHandler {
	return &JoinHandler{hub: hub}
}

func (h *LeaveHandler) Handle(ctx context.Context, client *Client, payload json.RawMessage) error {
	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}
	return h.hub.Broadcast(ctx, msg)
}
