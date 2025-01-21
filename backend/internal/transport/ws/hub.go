package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/jcserv/rivalslfg/internal/message"
	"github.com/jcserv/rivalslfg/internal/utils/log"
	"github.com/lxzan/gws"
)

type Message struct {
	GroupID  string             `json:"groupId"`
	PlayerID int                `json:"playerId"`
	Op       WebSocketEventType `json:"op"`
	Payload  interface{}        `json:"payload"`
}

type ClientInfo struct {
	GroupID  string `json:"groupId"`
	PlayerID int    `json:"playerId"`
}

type Hub struct {
	sync.RWMutex

	exchange message.Exchange

	// Map of group ID to set of client connections
	groups map[string]map[*Client]bool
	// Map of client to its current group ID
	clientGroups map[*Client]string
}

func NewHub(exchange message.Exchange) *Hub {
	return &Hub{
		exchange:     exchange,
		groups:       make(map[string]map[*Client]bool),
		clientGroups: make(map[*Client]string),
	}
}

func (h *Hub) Run(ctx context.Context) {
	pubsub := h.exchange.Subscribe(ctx)
	defer pubsub.Close()
	go h.handleRedisMessages(ctx, pubsub.Channel())

	<-ctx.Done()
	h.Lock()
	defer h.Unlock()

	// Clean up all connections when context is done
	for groupID, clients := range h.groups {
		for client := range clients {
			client.conn.NetConn().Close()
			delete(h.clientGroups, client)
		}
		delete(h.groups, groupID)
	}
}

func (h *Hub) RegisterClient(groupID string, client *Client) {
	h.Lock()
	defer h.Unlock()

	if h.groups[groupID] == nil {
		h.groups[groupID] = make(map[*Client]bool)
	}
	h.groups[groupID][client] = true
	h.clientGroups[client] = groupID
}

func (h *Hub) UnregisterClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	if groupID, ok := h.clientGroups[client]; ok {
		delete(h.clientGroups, client)
		if clients, exists := h.groups[groupID]; exists {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.groups, groupID)
			}
		}
	}
}

func (h *Hub) Broadcast(ctx context.Context, msg Message) error {
	h.RLock()
	defer h.RUnlock()

	if clients, exists := h.groups[msg.GroupID]; exists {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		for client := range clients {
			if err := client.conn.WriteMessage(gws.OpcodeText, msgBytes); err != nil {
				log.Warn(ctx, fmt.Sprintf("Error writing message to client: %v", err))
				h.UnregisterClient(client)
			}
		}
	}
	return nil
}

func (h *Hub) handleRedisMessages(ctx context.Context, ch <-chan *redis.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			event, err := message.UnmarshalJSON([]byte(msg.Payload))
			if err != nil {
				continue
			}
			h.RLock()

			// Broadcast event to clients in relevant group
			if clients, exists := h.groups[event.GroupID]; exists {
				wsMsg := Message{
					GroupID: event.GroupID,
					Op:      WebSocketEventType(event.Type),
					Payload: event.Payload,
				}

				for client := range clients {
					data, _ := json.Marshal(wsMsg)
					client.conn.WriteMessage(gws.OpcodeText, data)
				}
			}
			h.RUnlock()
		}
	}
}
