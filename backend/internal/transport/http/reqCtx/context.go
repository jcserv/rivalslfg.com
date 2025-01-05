package reqCtx

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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

func IsGroupOwner(ctx context.Context, groupID string) bool {
	claimGroupID := GetGroupID(ctx)
	return claimGroupID == groupID
}

func GetToken(ctx context.Context) string {
	info, ok := GetAuthInfo(ctx)
	if !ok {
		return ""
	}
	return info.Token
}

func ctxWithAuthInfo(ctx context.Context, info *AuthInfo) context.Context {
	return context.WithValue(ctx, authInfoKey, info)
}

func WithAuthInfo(r *http.Request, claims jwt.MapClaims) *http.Request {
	ctx := r.Context()

	playerID, groupID, token := 0, "", ""
	if playerIDVal, ok := claims["playerId"]; ok {
		playerID = utils.StringToInt(playerIDVal.(string))
	}
	if groupIDVal, ok := claims["groupId"]; ok {
		groupID = groupIDVal.(string)
	}
	if tokenVal, ok := claims["token"]; ok {
		token = tokenVal.(string)
	}

	info := &AuthInfo{
		PlayerID: playerID,
		GroupID:  groupID,
		Token:    token,
	}
	return r.WithContext(ctxWithAuthInfo(ctx, info))
}
