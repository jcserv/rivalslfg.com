package services

import (
	"context"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/message"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
)

type Group struct {
	repo      *repository.Queries
	publisher message.IPublisher
}

func NewGroup(repo *repository.Queries, publisher message.IPublisher) *Group {
	return &Group{
		repo:      repo,
		publisher: publisher,
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

	if group == nil {
		return nil, nil
	}

	if !isGroupOwner {
		group.Passcode = ""
	}

	return group, nil
}

func (s *Group) PatchGroup(ctx context.Context, arg repository.PatchGroupParams) (string, error) {
	result, err := s.repo.PatchGroup(ctx, arg)
	if err != nil {
		return "", err
	}

	if result == "404" {
		return "", NewError(http.StatusNotFound, "Group not found.", nil)
	}
	return result, nil
}

func (s *Group) DeleteGroup(ctx context.Context, id string) error {
	err := s.repo.DeleteGroup(ctx, id)
	if err != nil {
		return err
	}

	s.publisher.GroupDeleted(ctx, id, reqCtx.GetPlayerID(ctx))
	return nil
}
