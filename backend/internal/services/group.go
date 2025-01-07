package services

import (
	"context"

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

func (s *Group) GetGroupByID(ctx context.Context, id string, isGroupOwner bool) (*repository.GroupWithPlayers, error) {
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !isGroupOwner {
		group.Passcode = ""
	}

	return group, nil
}
