package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Right string

const (
	RightReadUser   Right = "user:read"
	RightUpdateUser Right = "user:update"
	RightDeleteUser Right = "user:delete"

	RightCreateGroup Right = "group:create"
	RightReadGroup   Right = "group:read"
	RightUpdateGroup Right = "group:update"
	RightDeleteGroup Right = "group:delete"

	RightJoinGroup  Right = "group:join"
	RightLeaveGroup Right = "group:leave"
)

// This is granted when a user joins a group, creates a group, or begins queuing up.
var BaseRights = []Right{
	RightReadUser,
	RightUpdateUser,
	RightDeleteUser,
	RightCreateGroup,
	RightReadGroup,
	RightJoinGroup,
}

var GroupMemberRights = []Right{
	RightLeaveGroup,
}

var GroupOwnerRights = []Right{
	RightUpdateGroup,
	RightDeleteGroup,
	RightLeaveGroup,
}

// TODO: How do we return the updated token after promoting a player?

func IsEqual(s string, r Right) bool {
	return string(r) == s
}

func HasOwnership(claims jwt.MapClaims, resourceType string, resourceID string) bool {
	switch resourceType {
	case "player":
		return claims["playerId"] == resourceID
	case "group":
		return claims["groupId"] == resourceID
	}
	return false
}
