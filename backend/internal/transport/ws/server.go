package ws

import (
	"context"
	"net/http"
)

type Server struct {
	hub     *Hub
	origins []string
}

func NewServer(hub *Hub, allowedOrigins []string) *Server {
	return &Server{
		hub:     hub,
		origins: allowedOrigins,
	}
}

func (s *Server) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(s.hub, w, r)
	})
}

func (s *Server) Start(ctx context.Context) {
	go s.hub.Run(ctx)
}
