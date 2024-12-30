package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/jcserv/rivalslfg/internal/transport/rest"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type Service struct {
	api *rest.API
	cfg *Configuration
}

func NewService() (*Service, error) {
	cfg, err := NewConfiguration()
	if err != nil {
		return nil, err
	}

	s := &Service{
		api: rest.NewAPI(),
		cfg: cfg,
	}
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

func (s *Service) StartHTTP(ctx context.Context) error {
	log.Info(ctx, fmt.Sprintf("Starting HTTP server on port %s", s.cfg.HTTPPort))
	r := s.api.RegisterRoutes()
	http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.HTTPPort), r)
	return nil
}
