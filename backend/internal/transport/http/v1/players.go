package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
	"github.com/jcserv/rivalslfg/internal/utils"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

// CreatePlayer: TODO
func (a *API) CreatePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.New("test error")
		httputil.InternalServerError(ctx, w, err)
		return

		// TODO: Generate token

		// ctx := r.Context()
		// token, err := a.authService.CreateAuth(ctx, "1")
		// if err != nil {
		// 	httputil.InternalServerError(ctx, w, err)
		// 	return
		// }

		// playerAuth, err := a.authService.ValidateToken(ctx, token)
		// if err != nil {
		// 	httputil.InternalServerError(ctx, w, err)
		// 	return
		// }

		// httputil.OK(w, map[string]any{
		// 	"token": string(token),
		// 	"playerAuth": map[string]string{
		// 		"playerId": playerAuth.PlayerID,
		// 		"token":    playerAuth.Token,
		// 		"lastSeen": playerAuth.LastSeen.String(),
		// 	},
		// })
	}
}

func (a *API) JoinGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input JoinGroup
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w, fmt.Errorf("unable to decode request body"))
			return
		}

		vars := mux.Vars(r)
		input.GroupID = vars["id"]
		input.PlayerID = reqCtx.GetPlayerID(ctx)

		params, err := input.Parse()
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w, err)
			return
		}

		playerID, err := a.playerService.JoinGroup(ctx, *params)
		if err != nil {
			if serviceErr, ok := err.(services.Error); ok {
				if serviceErr.Code() == http.StatusBadRequest {
					httputil.BadRequest(w, serviceErr)
					return
				}
				if serviceErr.Code() == http.StatusNotFound {
					httputil.NotFound(w)
					return
				}
				if serviceErr.Code() == http.StatusForbidden {
					httputil.Forbidden(w)
					return
				}
			}
			httputil.InternalServerError(ctx, w, err)
			return
		}

		httputil.EmbedTokenInResponse(ctx, w, &reqCtx.AuthInfo{
			PlayerID: int(playerID),
			GroupID:  input.GroupID,
		}, auth.GroupMemberRights)

		httputil.OK(w, nil)
	}
}

func (a *API) RemovePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		input := RemovePlayer{
			GroupID:          vars["id"],
			PlayerToRemoveID: utils.StringToInt(vars["playerId"]),
		}

		params, err := input.Parse()
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w, err)
			return
		}

		requesterID := reqCtx.GetPlayerID(ctx)
		// Users cannot remove players from a group you are not in
		if requesterID != 0 && !reqCtx.IsGroupMember(ctx, input.GroupID) {
			httputil.Forbidden(w)
			return
		}

		// Users can only remove yourself if you're not the owner of the group
		if requesterID == 0 || (requesterID != input.PlayerToRemoveID && !reqCtx.IsGroupOwner(ctx, input.GroupID)) {
			httputil.Forbidden(w)
			return
		}

		status, err := a.playerService.RemovePlayer(ctx, *params)
		if err != nil {
			if serviceErr, ok := err.(services.Error); ok {
				switch serviceErr.Code() {
				case http.StatusBadRequest:
					httputil.BadRequest(w, serviceErr)
					return
				case http.StatusNotFound:
					httputil.NotFound(w)
					return
				}
			}
			httputil.InternalServerError(ctx, w, err)
			return
		}

		httputil.EmbedTokenInResponse(ctx, w, &reqCtx.AuthInfo{
			PlayerID: 0,
			GroupID:  "",
		}, []auth.Right{})

		if status == "204" {
			httputil.NoContent(w)
			return
		}

		httputil.OK(w, nil)
	}
}
