package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/lxzan/gws"
)

const (
	OpGroupChat = iota + 1
	OpGroupJoin
	OpGroupLeave
	OpGroupPromotion
)

type Message struct {
	GroupID string      `json:"groupId"`
	Op      int         `json:"op"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	sync.RWMutex
	// Map of group ID to set of client connections
	groups map[string]map[*Client]bool
	// Map of client to its current group ID
	clientGroups map[*Client]string
}

func NewHub() *Hub {
	return &Hub{
		groups:       make(map[string]map[*Client]bool),
		clientGroups: make(map[*Client]string),
	}
}

func (h *Hub) Run(ctx context.Context) {
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
