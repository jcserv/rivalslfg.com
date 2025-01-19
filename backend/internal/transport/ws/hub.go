package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/jcserv/rivalslfg/internal/message"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
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
	clientGroups map[*Client]*ClientInfo
}

func NewHub(exchange message.Exchange) *Hub {
	return &Hub{
		exchange:     exchange,
		groups:       make(map[string]map[*Client]bool),
		clientGroups: make(map[*Client]*ClientInfo),
	}
}

func (h *Hub) Run(ctx context.Context) {
	pubsub, err := h.exchange.Subscribe(ctx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error subscribing to group updates: %v", err))
	} else {
		log.Info(ctx, "Subscribed to group updates channel")
		defer pubsub.Close()
		go h.handleRedisMessages(ctx, pubsub.Channel())
	}

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

func (h *Hub) RegisterClient(authInfo *reqCtx.AuthInfo, client *Client) {
	h.Lock()
	defer h.Unlock()

	if h.groups[authInfo.GroupID] == nil {
		h.groups[authInfo.GroupID] = make(map[*Client]bool)
	}
	h.groups[authInfo.GroupID][client] = true
	h.clientGroups[client] = &ClientInfo{
		GroupID:  authInfo.GroupID,
		PlayerID: authInfo.PlayerID,
	}
}

func (h *Hub) UnregisterClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	if clientInfo, ok := h.clientGroups[client]; ok {
		delete(h.clientGroups, client)
		if clients, exists := h.groups[clientInfo.GroupID]; exists {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.groups, clientInfo.GroupID)
			}
		}
	}
}

func (h *Hub) Broadcast(msg Message) error {
	h.RLock()
	defer h.RUnlock()

	clients, exists := h.groups[msg.GroupID]
	if !exists {
		return nil
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for client := range clients {
		client.conn.WriteMessage(gws.OpcodeText, msgBytes)
	}
	return nil
}

func (h *Hub) handleRedisMessages(ctx context.Context, ch <-chan *redis.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			log.Info(ctx, fmt.Sprintf("Received message from Redis: %s", msg.Payload))

			event, err := message.UnmarshalJSON([]byte(msg.Payload))
			if err != nil {
				continue
			}
			h.RLock()

			// Broadcast event to clients in relevant group
			if clients, exists := h.groups[event.GroupID]; exists {
				log.Info(ctx, fmt.Sprintf("Broadcasting message to %d clients", len(clients)))
				wsMsg := Message{
					GroupID: event.GroupID,
					Op:      WebSocketEventType(event.Type),
					Payload: event.Payload,
				}

				for client := range clients {
					info, exists := h.clientGroups[client]
					if !exists {
						continue
					}

					if info.PlayerID != event.PlayerID {
						log.Info(ctx, fmt.Sprintf("Sending message to client %d", info.PlayerID))
						data, _ := json.Marshal(wsMsg)
						client.conn.WriteMessage(gws.OpcodeText, data)
					}
				}
			}
			h.RUnlock()
		}
	}
}
