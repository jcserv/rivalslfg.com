package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type Group struct {
	repo *repository.Queries
}

func NewGroup(repo *repository.Queries) *Group {
	return &Group{
		repo: repo,
	}
}

func (s *Group) UpsertGroup(ctx context.Context, arg repository.UpsertGroupParams) (string, error) {
	groupID, err := s.repo.UpsertGroup(ctx, arg)
	if err != nil {
		return "", err
	}
	return groupID, nil
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

type JoinGroupArgs struct {
	CheckCanJoinGroup repository.CheckCanJoinGroupParams
	JoinGroup         repository.JoinGroupParams
}

func (s *Group) JoinGroup(ctx context.Context, arg JoinGroupArgs) error {
	// TODO: Acquire lock on group
	// TODO: Check against requirements of group
	status, err := s.repo.CheckCanJoinGroup(ctx, arg.CheckCanJoinGroup)
	if err != nil {
		return err
	}

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

func (s *Group) RemovePlayerFromGroup(ctx context.Context, arg repository.RemovePlayerFromGroupParams) error {
	out, err := s.repo.RemovePlayerFromGroup(ctx, arg)
	if err != nil {
		return err
	}

	switch out.Status {
	case 200:
		if !(reqCtx.IsGroupOwner(ctx, arg.ID)) {
			return nil
		}
		return s.PromoteOwnerOrDeleteGroup(ctx, repository.PromoteOwnerOrDeleteGroupParams{
			ID:               arg.ID,
			RemainingPlayers: out.RemainingPlayers,
		})
	case 403:
		return NewError(http.StatusForbidden, "Only group owners can remove other players", nil)
	case 404:
		return NewError(http.StatusNotFound, "Group or player not found", nil)
	default:
		log.Debug(ctx, fmt.Sprintf("Unknown error: %v", out.Status))
		return NewError(http.StatusInternalServerError, "Unknown error", nil)
	}
}

func (s *Group) PromoteOwnerOrDeleteGroup(ctx context.Context, arg repository.PromoteOwnerOrDeleteGroupParams) error {
	out, err := s.repo.PromoteOwnerOrDeleteGroup(ctx, arg)
	if err != nil {
		return err
	}

	switch out.(type) {
	case string:
		switch out.(string) {
		case "404":
			return NewError(http.StatusNotFound, "Group not found", nil)
		default:
			return nil
		}
	default:
		log.Debug(ctx, fmt.Sprintf("Unknown error: %v", out))
		return NewError(http.StatusInternalServerError, "Unknown error", nil)
	}
}
