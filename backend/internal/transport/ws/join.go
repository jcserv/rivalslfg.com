package ws

import (
	"context"
	"encoding/json"

	"github.com/jcserv/rivalslfg/internal/utils/log"
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
	log.Info(ctx, "Broadcasting join message")
	return h.hub.Broadcast(msg)
}
