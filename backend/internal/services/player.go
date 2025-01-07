package services

import (
	"context"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type Player struct {
	repo *repository.Queries
}

func NewPlayer(repo *repository.Queries) *Player {
	return &Player{
		repo: repo,
	}
}

func (s *Player) JoinGroup(ctx context.Context, arg repository.JoinGroupParams) (int32, error) {
	result, err := s.repo.JoinGroup(ctx, arg)
	if err != nil {
		return 0, err
	}

	switch result.Status {
	case "200":
		return result.PlayerID, nil
	case "400a":
		return 0, NewError(http.StatusBadRequest, "Player is already in a group.", nil)
	case "404":
		return 0, NewError(http.StatusNotFound, "Group not found.", nil)
	case "403":
		return 0, NewError(http.StatusForbidden, "Access denied.", nil)
	case "400e":
		return 0, NewError(http.StatusBadRequest, "Group requirements not met.", nil)
	default:
		return 0, NewError(http.StatusInternalServerError, "An unexpected error occurred.", nil)
	}
}
