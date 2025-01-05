package v1

import (
	"errors"
	"net/http"

	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
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
		return

		// ctx := r.Context()

		// var input JoinGroup
		// err := json.NewDecoder(r.Body).Decode(&input)
		// if err != nil {
		// 	log.Debug(ctx, err.Error())
		// 	httputil.BadRequest(w)
		// 	return
		// }

		// vars := mux.Vars(r)
		// groupID := vars["id"]
		// input.GroupID = groupID

		// params, err := input.Parse()
		// if err != nil {
		// 	log.Debug(ctx, err.Error())
		// 	httputil.BadRequest(w)
		// 	return
		// }

		// err = a.groupService.JoinGroup(ctx, *params)
		// if err != nil {
		// 	if serviceErr, ok := err.(services.Error); ok {
		// 		if serviceErr.Code() == http.StatusNotFound {
		// 			httputil.NotFound(w)
		// 			return
		// 		}
		// 		if serviceErr.Code() == http.StatusForbidden {
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
