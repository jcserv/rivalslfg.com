package ws

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	hub     *Hub
	port    string
	origins []string
}

func NewServer(port string, allowedOrigins []string) *Server {
	return &Server{
		hub:     NewHub(),
		port:    port,
		origins: allowedOrigins,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go s.hub.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(s.hub, w, r)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down WebSocket server: %v\n", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("WebSocket server error: %v", err)
	}

	return nil
}

func (s *Server) Broadcast(msg Message) error {
	return s.hub.Broadcast(msg)
}
