package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/lxzan/gws"
)

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

type Client struct {
	hub           *Hub
	conn          *gws.Conn
	eventHandlers map[WebSocketEventType]EventHandler
}

func NewClient(hub *Hub, conn *gws.Conn) *Client {
	client := &Client{
		hub:           hub,
		conn:          conn,
		eventHandlers: make(map[WebSocketEventType]EventHandler),
	}

	// Register default handlers
	client.eventHandlers[OpGroupChat] = NewChatHandler(hub)

	return client
}

type ClientHandler struct {
	hub    *Hub
	client *Client
}

func (h *ClientHandler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	h.client = NewClient(h.hub, socket)
}

func (h *ClientHandler) OnClose(socket *gws.Conn, err error) {
	if err != nil {
	}
	if h.client != nil {
		h.hub.UnregisterClient(h.client)
	}
}

func (h *ClientHandler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}

func (h *ClientHandler) OnPong(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (h *ClientHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	var msg Message
	if err := json.Unmarshal(message.Bytes(), &msg); err != nil {
		return
	}

	if handler, exists := h.client.eventHandlers[msg.Op]; exists {
		if err := handler.Handle(h.client, message.Bytes()); err != nil {
			return
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("groupId")
	if groupID == "" {
		return
	}

	client := &Client{
		hub: hub,
	}

	handler := &ClientHandler{
		hub:    hub,
		client: client,
	}

	loggingHandler := NewLoggingMiddleware(handler)

	upgrader := gws.NewUpgrader(loggingHandler, &gws.ServerOption{
		ParallelEnabled:   true,
		Recovery:          gws.Recovery,
		PermessageDeflate: gws.PermessageDeflate{Enabled: true},
	})

	conn, err := upgrader.Upgrade(w, r)
	if err != nil {
		return
	}

	client.conn = conn
	hub.RegisterClient(groupID, client)

	go func() {
		conn.ReadLoop() // Blocking prevents the context from being GC
	}()
}
