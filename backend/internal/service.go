package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/transport/rest"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type Service struct {
	api          *rest.API
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

	s.api = rest.NewAPI(services.NewGroupService(repository.New(conn)))
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
	http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.HTTPPort), r)
	return nil
}
