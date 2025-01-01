package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
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
			log.Error(ctx, err.Error())
			httputil.InternalServerError(ctx, w, err)
			return
		}

		httputil.OK(w, map[string]string{
			"id": groupID,
		})
	}
}

func (a *API) ReadGroup() http.HandlerFunc {
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
			log.Error(ctx, err.Error())
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

// TODO: Add pagination
func (a *API) ReadGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		groups, err := a.groupService.GetGroups(ctx)
		if err != nil {
			log.Error(ctx, err.Error())
			httputil.InternalServerError(ctx, w, err)
			return
		}

		httputil.OK(w, groups)
	}
}

func (a *API) DeleteGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (a *API) JoinGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
