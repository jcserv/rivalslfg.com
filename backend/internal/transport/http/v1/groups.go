package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
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
			httputil.InternalServerError(ctx, w, err)
			return
		}

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
			httputil.BadRequest(w)
			return
		}

		args, err := Parse(queryParams)
		if err != nil {
			httputil.BadRequest(w)
			return
		}

		groups, err := a.groupService.GetGroups(ctx, *args)
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			return
		}

		httputil.OK(w, groups)
	}
}

func (a *API) DeleteGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errors.New("test error")
		httputil.InternalServerError(ctx, w, err)
		return
	}
}

func (a *API) JoinGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input JoinGroup
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w)
			return
		}

		vars := mux.Vars(r)
		groupID := vars["id"]
		input.GroupID = groupID

		params, err := input.Parse()
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w)
			return
		}

		err = a.groupService.JoinGroup(ctx, *params)
		if err != nil {
			if serviceErr, ok := err.(services.Error); ok {
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
		httputil.OK(w, nil)
	}
}

func (a *API) RemovePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		groupID := vars["id"]

		var input RemovePlayer
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w)
			return
		}

		input.GroupID = groupID
		params, err := input.Parse()
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w)
			return
		}

		err = a.groupService.RemovePlayerFromGroup(ctx, *params)
		if err != nil {
			if serviceErr, ok := err.(services.Error); ok {
				switch serviceErr.Code() {
				case http.StatusNotFound:
					httputil.NotFound(w)
					return
				case http.StatusForbidden:
					httputil.Forbidden(w)
					return
				}
			}
			log.Error(ctx, err.Error())
			httputil.InternalServerError(ctx, w, err)
			return
		}

		httputil.OK(w, nil)
	}
}
