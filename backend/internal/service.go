package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/services"
	_http "github.com/jcserv/rivalslfg/internal/transport/http"
	v1 "github.com/jcserv/rivalslfg/internal/transport/http/v1"
	"github.com/jcserv/rivalslfg/internal/transport/ws"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type Service struct {
	api *_http.API
	cfg *Configuration
}

func NewService() (*Service, error) {
	cfg, err := NewConfiguration()
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	s := &Service{
		cfg: cfg,
	}

	conn, err := s.ConnectDB(context.Background())
	if err != nil {
		return nil, err
	}

	// client, err := s.ConnectCache(context.Background())
	// if err != nil {
	// 	return nil, err
	// }

	repo := repository.New(conn)
	// store := store.New(client)

	s.api = _http.NewAPI(
		&v1.Dependencies{
			GroupService:  services.NewGroup(repo),
			PlayerService: services.NewPlayer(repo),
		},
	)
	return s, nil
}

func (s *Service) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		s.StartHTTP(ctx)
	}(ctx)
	wg.Add(1)

	wg.Wait()
	return nil
}

func (s *Service) ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(s.cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	log.Info(ctx, "Connection established to database")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return pool, nil
}

func (s *Service) ConnectCache(ctx context.Context) (*redis.Client, error) {
	config, err := redis.ParseURL(s.cfg.CacheURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(config)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	log.Info(ctx, "Connection established to cache")
	return client, nil
}

func (s *Service) StartHTTP(ctx context.Context) error {
	log.Info(ctx, fmt.Sprintf("Starting HTTP server on port %s", s.cfg.HTTPPort))

	wsServer := ws.NewServer([]string{os.Getenv("ORIGIN_ALLOWED")})
	wsServer.Start(ctx)

	mainMux := http.NewServeMux()
	r := s.api.RegisterRoutes()
	mainMux.Handle("/", r)
	wsServer.RegisterHandlers(mainMux)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})

	exposedHeaders := handlers.ExposedHeaders([]string{"X-Requested-With", "Content-Type", "X-Total-Count", "X-Token"})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.HTTPPort),
		Handler: handlers.CORS(origins, headers, methods, exposedHeaders)(mainMux),
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(ctx, fmt.Sprintf("Error shutting down HTTP server: %v", err))
		}
	}()

	srv.ListenAndServe()
	return nil
}
