package ws

import (
	"context"
	"encoding/json"

	"github.com/lxzan/gws"
	"go.uber.org/zap"

	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type ClientMiddleware interface {
	OnConnect(socket *gws.Conn)
	OnMessage(message *Message) error
	OnClose(socket *gws.Conn, err error)
}

type LoggingMiddleware struct {
	logger *zap.Logger
	next   *ClientHandler
}

func NewLoggingMiddleware(next *ClientHandler) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: log.GetLogger(context.Background()),
		next:   next,
	}
}

func (m *LoggingMiddleware) OnOpen(socket *gws.Conn) {
	m.next.OnOpen(socket)
}

func (m *LoggingMiddleware) OnClose(socket *gws.Conn, err error) {
	m.next.OnClose(socket, err)
}

func (m *LoggingMiddleware) OnPing(socket *gws.Conn, payload []byte) {
	m.next.OnPing(socket, payload)
}

func (m *LoggingMiddleware) OnPong(socket *gws.Conn, payload []byte) {
	m.next.OnPong(socket, payload)
}

func (m *LoggingMiddleware) OnMessage(socket *gws.Conn, message *gws.Message) {
	var msg Message
	if err := json.Unmarshal(message.Bytes(), &msg); err != nil {
		m.logger.Error("Failed to parse WebSocket message",
			zap.Error(err),
			zap.String("remote_addr", socket.RemoteAddr().String()),
		)
		return
	}

	m.logger.Info("WebSocket message received",
		zap.String("remote_addr", socket.RemoteAddr().String()),
		zap.String("group_id", msg.GroupID),
		zap.Int("op", int(msg.Op)),
		zap.Any("payload", msg.Payload),
	)

	m.next.OnMessage(socket, message)
}
