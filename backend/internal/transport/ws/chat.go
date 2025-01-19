package ws

import (
	"context"
	"encoding/json"
)

type ChatPayload struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
}

type ChatHandler struct {
	hub *Hub
}

func NewChatHandler(hub *Hub) *ChatHandler {
	return &ChatHandler{hub: hub}
}

func (h *ChatHandler) Handle(ctx context.Context, client *Client, payload json.RawMessage) error {
	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return err
	}
	return h.hub.Broadcast(msg)
}
