package services

import (
	"context"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type GroupService struct {
	repo *repository.Queries
}

func NewGroupService(repo *repository.Queries) *GroupService {
	return &GroupService{
		repo: repo,
	}
}

func (s *GroupService) UpsertGroup(ctx context.Context, arg repository.UpsertGroupParams) (string, error) {
	groupID, err := s.repo.UpsertGroup(ctx, arg)
	if err != nil {
		return "", err
	}
	return groupID, nil
}

func (s *GroupService) GetGroups(ctx context.Context) ([]repository.GroupWithPlayers, error) {
	groups, err := s.repo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *GroupService) GetGroupByID(ctx context.Context, id string) (*repository.GroupWithPlayers, error) {
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) JoinGroup(ctx context.Context, arg repository.JoinGroupParams) error {
	// TODO: Acquire lock on group

	result, err := s.repo.JoinGroup(ctx, arg)
	if err != nil {
		return err
	}
	switch result {
	case 200:
		return nil
	case 403:
		return NewError(http.StatusForbidden, "Passcode does not match", nil)
	case 404:
		return NewError(http.StatusNotFound, "Group not found", nil)
	default:
		return NewError(http.StatusInternalServerError, "Unknown error", nil)
	}

	// TODO: Emit message to notify players
}
