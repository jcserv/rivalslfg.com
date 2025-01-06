package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

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

		// TODO: Change player.roles to single role
		err = a.groupService.JoinGroup(ctx, *params)
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
				httputil.InternalServerError(ctx, w, err)
			}
			return
		}

		// httputil.EmbedTokenInResponse(ctx, w, &reqCtx.AuthInfo{
		// 	PlayerID: int(result.PlayerID),
		// 	GroupID:  result.GroupID,
		// }, auth.GroupMemberRights)
		httputil.OK(w, nil)
	}
}

func (a *API) RemovePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
		// ctx := r.Context()

		// vars := mux.Vars(r)
		// input := RemovePlayer{
		// 	GroupID:          vars["id"],
		// 	RequesterID:      reqCtx.GetPlayerID(ctx),
		// 	PlayerToRemoveID: vars["playerId"],
		// }

		// params, err := input.Parse()
		// if err != nil {
		// 	log.Debug(ctx, err.Error())
		// 	httputil.BadRequest(w)
		// 	return
		// }

		// err = a.groupService.RemovePlayerFromGroup(ctx, *params)
		// if err != nil {
		// 	if serviceErr, ok := err.(services.Error); ok {
		// 		switch serviceErr.Code() {
		// 		case http.StatusNotFound:
		// 			httputil.NotFound(w)
		// 			return
		// 		case http.StatusForbidden:
		// 			httputil.Forbidden(w)
		// 			return
		// 		}
		// 	}
		// 	httputil.InternalServerError(ctx, w, err)
		// 	return
		// }

		// httputil.OK(w, nil)
	}
}
