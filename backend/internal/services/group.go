package services

import (
	"context"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type Group struct {
	repo *repository.Queries
}

func NewGroup(repo *repository.Queries) *Group {
	return &Group{
		repo: repo,
	}
}

func (s *Group) CreateGroup(ctx context.Context, arg repository.CreateGroupParams) (repository.CreateGroupRow, error) {
	result, err := s.repo.CreateGroup(ctx, arg)
	if err != nil {
		return repository.CreateGroupRow{}, err
	}
	return result, nil
}

func (s *Group) GetGroups(ctx context.Context, arg repository.GetGroupsParams) ([]repository.GroupWithPlayers, int32, error) {
	result, err := s.repo.GetGroups(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	return result.Groups, result.TotalCount, nil
}

func (s *Group) GetGroupByID(ctx context.Context, id string) (*repository.GroupWithPlayers, error) {
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *Group) JoinGroup(ctx context.Context, arg repository.JoinGroupParams) error {
	status, err := s.repo.JoinGroup(ctx, arg)
	if err != nil {
		return err
	}

	switch status {
	case "200":
		return nil
	case "400a":
		return NewError(http.StatusBadRequest, "Player is already in a group.", nil)
	case "404":
		return NewError(http.StatusNotFound, "Group not found.", nil)
	case "403":
		return NewError(http.StatusForbidden, "Access denied.", nil)
	case "400e":
		return NewError(http.StatusBadRequest, "Group requirements not met.", nil)
	default:
		return NewError(http.StatusInternalServerError, "An unexpected error occurred.", nil)
	}
}
