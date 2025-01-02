package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v5"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/services"
	_http "github.com/jcserv/rivalslfg/internal/transport/http"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type Service struct {
	api          *_http.API
	cfg          *Configuration
	groupService *services.GroupService
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

	repo := repository.New(conn)
	s.api = _http.NewAPI(services.NewGroupService(repo))
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

	wg.Wait()
	return nil
}

func (s *Service) ConnectDB(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), s.cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	log.Info(ctx, "Connection established to database")

	customDataTypes := []string{"RankName", "RankID", "Rank"}
	for _, typeName := range customDataTypes {
		dataType, _ := conn.LoadType(ctx, typeName)
		conn.TypeMap().RegisterType(dataType)
	}

	return conn, nil
}

func (s *Service) StartHTTP(ctx context.Context) error {
	log.Info(ctx, fmt.Sprintf("Starting HTTP server on port %s", s.cfg.HTTPPort))
	r := s.api.RegisterRoutes()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})

	exposedHeaders := handlers.ExposedHeaders([]string{"X-Requested-With", "Content-Type", "X-Total-Count"})

	http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.HTTPPort), handlers.CORS(origins, headers, methods, exposedHeaders)(r))
	return nil
}
