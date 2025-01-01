package services

import (
	"context"

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

// func (s *GroupService) CreateGroupWithOwner(ctx context.Context, arg repository.CreateGroupWithOwnerParams) (string, error) {
// 	groupID, err := s.repo.CreateGroupWithOwner(ctx, arg)
// 	if err != nil {
// 		return "", err
// 	}
// 	return groupID, nil
// }

func (s *GroupService) GetGroups(ctx context.Context) ([]repository.GroupWithPlayers, error) {
	groups, err := s.repo.FindAllGroups(ctx)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *GroupService) GetGroupByID(ctx context.Context, id string) (*repository.GroupWithPlayers, error) {
	group, err := s.repo.FindGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return group, nil
}
