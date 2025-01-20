package ws

import (
	"context"
	"encoding/json"
)

type DeleteGroupPayload struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
}

type DeleteGroupHandler struct {
	hub *Hub
}

func NewDeleteGroupHandler(hub *Hub) *DeleteGroupHandler {
	return &DeleteGroupHandler{hub: hub}
}

func (h *DeleteGroupHandler) Handle(ctx context.Context, client *Client, payload json.RawMessage) error {
	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}
	return h.hub.Broadcast(ctx, msg)
}
