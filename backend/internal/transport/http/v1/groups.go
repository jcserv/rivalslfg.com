package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

func (a *API) UpsertGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input UpsertGroup
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w)
			return
		}

		params, err := input.Parse()
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w)
			return
		}

		groupID, err := a.groupService.UpsertGroup(ctx, *params)
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			return
		}

		authInfo := reqCtx.GetAuthInfoOrDefault(ctx, &reqCtx.AuthInfo{
			PlayerID: "1", // TODO: This should be generated server-side
			GroupID:  groupID,
		})
		httputil.EmbedTokenInResponse(ctx, w, authInfo, auth.GroupOwnerRights)
		httputil.OK(w, map[string]string{
			"id": groupID,
		})
	}
}

func (a *API) GetGroupByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		groupID := vars["id"]
		if groupID == "" {
			httputil.BadRequest(w)
			return
		}

		group, err := a.groupService.GetGroupByID(ctx, groupID)
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			return
		}

		if group == nil {
			httputil.NotFound(w)
			return
		}

		httputil.OK(w, group)
	}
}

func (a *API) GetGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		queryParams, err := httputil.ParseQueryParams(r)
		if err != nil {
			log.Debug(ctx, fmt.Sprintf("error parsing query params: %v", err))
			httputil.BadRequest(w)
			return
		}

		args, err := Parse(queryParams)
		if err != nil {
			log.Debug(ctx, fmt.Sprintf("error parsing query params: %v", err))
			httputil.BadRequest(w)
			return
		}

		groups, totalCount, err := a.groupService.GetGroups(ctx, *args)
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			return
		}

		if queryParams.PaginateBy.Count {
			w.Header().Set("X-Total-Count", strconv.FormatInt(int64(totalCount), 10))
		}

		httputil.OK(w, groups)
	}
}

func (a *API) DeleteGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.New("test error")
		httputil.InternalServerError(ctx, w, err)

		// TODO: Generate token with groupId = "", and remove group owner rights
		return
	}
}
