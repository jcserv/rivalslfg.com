package services

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type PlayerService struct {
	repo *repository.Queries
}

func NewPlayerService(repo *repository.Queries) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (s *PlayerService) UpsertPlayer(ctx context.Context, arg repository.UpsertPlayerParams) (int32, error) {
	return s.repo.UpsertPlayer(ctx, arg)
}

func (s *PlayerService) FindPlayer(ctx context.Context, id int32, name string) (*repository.Player, error) {
	player, err := s.repo.FindPlayer(ctx, id, name)
	if err != nil {
		return nil, err
	}
	return player, nil
}
