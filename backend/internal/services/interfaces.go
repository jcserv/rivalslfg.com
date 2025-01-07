package services

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type IAuth interface {
	CreateAuth(ctx context.Context, playerID string) (Token, error)
	ValidateToken(ctx context.Context, token Token) (*PlayerAuth, error)
}

type IGroup interface {
	CreateGroup(ctx context.Context, arg repository.CreateGroupParams) (repository.CreateGroupRow, error)
	GetGroups(ctx context.Context, arg repository.GetGroupsParams) ([]repository.GroupWithPlayers, int32, error)
	GetGroupByID(ctx context.Context, id string) (*repository.GroupWithPlayers, error)
	JoinGroup(ctx context.Context, arg repository.JoinGroupParams) (int32, error)
	// RemovePlayerFromGroup(ctx context.Context, arg repository.RemovePlayerFromGroupParams) error
	// PromoteOwnerOrDeleteGroup(ctx context.Context, arg repository.PromoteOwnerOrDeleteGroupParams) error
}
