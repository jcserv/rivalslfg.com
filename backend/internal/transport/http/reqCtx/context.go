package reqCtx

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/utils"
)

type AuthInfo struct {
	PlayerID int
	GroupID  string
	Token    string
}

type contextKey string

const authInfoKey contextKey = "authInfo"

func GetAuthInfo(ctx context.Context) (*AuthInfo, bool) {
	val := ctx.Value(authInfoKey)
	if val == nil {
		return nil, false
	}

	info, ok := val.(*AuthInfo)
	if !ok || info == nil {
		return nil, false
	}
	return info, ok
}

func GetAuthInfoOrDefault(ctx context.Context, fallback *AuthInfo) *AuthInfo {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return fallback
	}
	return info
}

func GetPlayerID(ctx context.Context) int {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return 0
	}
	return info.PlayerID
}

func GetGroupID(ctx context.Context) string {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return ""
	}
	return info.GroupID
}

func GetToken(ctx context.Context) string {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return ""
	}
	return info.Token
}

func IsGroupOwner(ctx context.Context, groupID string) bool {
	if GetGroupID(ctx) != groupID {
		return false
	}

	token := GetToken(ctx)
	if token == "" {
		return false
	}
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return false
	}
	return auth.HasRight(claims, auth.RightDeleteGroup)
}

func IsGroupMember(ctx context.Context, groupID string) bool {
	if GetGroupID(ctx) != groupID {
		return false
	}

	token := GetToken(ctx)
	if token == "" {
		return false
	}
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return false
	}
	return auth.HasRight(claims, auth.RightLeaveGroup)
}

func ctxWithAuthInfo(ctx context.Context, info *AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoKey, info)
}

func Init(r *http.Request, claims jwt.MapClaims, token string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "request_id", uuid.New().String())

	playerID, groupID := 0, ""
	if playerIDVal, ok := claims["playerId"]; ok {
		playerID = utils.StringToInt(playerIDVal.(string))
	}
	if groupIDVal, ok := claims["groupId"]; ok {
		groupID = groupIDVal.(string)
	}

	info := &AuthInfo{
		PlayerID: playerID,
		GroupID:  groupID,
		Token:    token,
	}
	return r.WithContext(ctxWithAuthInfo(ctx, info))
}

func WithAuthInfo(r *http.Request, info *AuthInfo) *http.Request {
	ctx := r.Context()
	return r.WithContext(ctxWithAuthInfo(ctx, info))
}
