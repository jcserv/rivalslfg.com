package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/utils/log"
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

type JoinGroupArgs struct {
	CheckCanJoinGroup repository.CheckCanJoinGroupParams
	JoinGroup         repository.JoinGroupParams
}

func (s *GroupService) JoinGroup(ctx context.Context, arg JoinGroupArgs) error {
	// TODO: Acquire lock on group
	// TODO: Check against requirements of group
	status, err := s.repo.CheckCanJoinGroup(ctx, arg.CheckCanJoinGroup)
	if err != nil {
		return err
	}

	log.Debug(ctx, fmt.Sprintf("status %d", status))

	switch status {
	case http.StatusOK: // User already in group
		return nil
	case http.StatusAccepted:
		if err := s.repo.JoinGroup(ctx, arg.JoinGroup); err != nil {
			return NewError(http.StatusInternalServerError, "Failed to add player.", err)
		}
		// TODO: Emit message to notify players
		return nil
	case http.StatusForbidden:
		return NewError(http.StatusForbidden, "Passcode does not match.", nil)
	case http.StatusNotFound:
		return NewError(http.StatusNotFound, "Group not found.", nil)
	default:
		return NewError(http.StatusInternalServerError, "An unexpected error occurred.", nil)
	}

}
