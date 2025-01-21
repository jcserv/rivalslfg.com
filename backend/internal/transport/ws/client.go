package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jcserv/rivalslfg/internal/utils/log"
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

	client.eventHandlers[OpGroupChat] = NewDefaultHandler(hub)
	client.eventHandlers[OpGroupJoin] = NewDefaultHandler(hub)
	client.eventHandlers[OpGroupLeave] = NewDefaultHandler(hub)
	client.eventHandlers[OpGroupDelete] = NewDefaultHandler(hub)

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

func (h *ClientHandler) OnClose(socket *gws.Conn, _ error) {
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

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", uuid.New().String())

	var msg Message
	if err := json.Unmarshal(message.Bytes(), &msg); err != nil {
		return
	}

	if _, ok := h.hub.clientGroups[h.client]; !ok {
		log.Warn(ctx, "Received message from unregistered client")
		return
	}

	if handler, exists := h.client.eventHandlers[msg.Op]; exists {
		if err := handler.Handle(ctx, h.client, message.Bytes()); err != nil {
			return
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	client := &Client{
		hub: hub,
	}

	handler := &ClientHandler{
		hub:    hub,
		client: client,
	}

	loggingHandler := NewLoggingMiddleware(handler)

	upgrader := gws.NewUpgrader(loggingHandler, &gws.ServerOption{
		ParallelEnabled: true,
		Recovery:        gws.Recovery,
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})

	conn, err := upgrader.Upgrade(w, r)
	if err != nil {
		return
	}

	client.conn = conn

	groupId := r.URL.Query().Get("groupId")
	if groupId == "" {
		return
	}

	hub.RegisterClient(groupId, client)

	go func() {
		conn.ReadLoop() // Blocking prevents the context from being GC
	}()
}
