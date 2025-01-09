package services

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/repository"
)

type IGroup interface {
	CreateGroup(ctx context.Context, arg repository.CreateGroupParams) (repository.CreateGroupRow, error)
	GetGroups(ctx context.Context, arg repository.GetGroupsParams) ([]repository.GroupWithPlayers, int32, error)
	GetGroupByID(ctx context.Context, id string, isGroupOwner bool) (*repository.GroupWithPlayers, error)
}

type IPlayer interface {
	JoinGroup(ctx context.Context, arg repository.JoinGroupParams) (int32, error)
	RemovePlayer(ctx context.Context, arg repository.RemovePlayerParams) (string, error)
}
