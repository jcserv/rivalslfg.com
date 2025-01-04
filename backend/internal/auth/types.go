package auth

type Right string

const (
	RightReadUser   Right = "user:read"
	RightUpdateUser Right = "user:update"
	RightDeleteUser Right = "user:delete"

	RightCreateGroup Right = "group:create"
	RightReadGroup   Right = "group:read"
	RightUpdateGroup Right = "group:update"
	RightDeleteGroup Right = "group:delete"
	RightLeaveGroup  Right = "group:leave"
)

func IsEqual(s string, r Right) bool {
	return string(r) == s
}
