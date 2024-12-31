package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

const (
	APIV1URLPath = "/api/v1/"
	groups       = APIV1URLPath + "groups"
	group        = APIV1URLPath + "groups/{id}"
	joinGroup    = group + "/join"
)

type API struct {
	groupService *services.GroupService
}

func NewAPI(groupService *services.GroupService) *API {
	return &API{
		groupService: groupService,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(groups, a.CreateGroup()).Methods(http.MethodPost)
	r.HandleFunc(group, a.ReadGroup()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.ReadGroups()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.UpdateGroup()).Methods(http.MethodPut)
	r.HandleFunc(groups, a.DeleteGroup()).Methods(http.MethodDelete)
	r.HandleFunc(joinGroup, a.JoinGroup()).Methods(http.MethodPost)
}

func (a *API) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input CreateGroup
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			httputil.BadRequest(w)
			return
		}

		if err := input.Validate(); err != nil {
			httputil.BadRequest(w)
			return
		}

		groupID, err := a.groupService.CreateGroupWithOwner(ctx, input.ToParams())
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
			return
		}

		httputil.OK(w, groupID)
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
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
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
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
			return
		}

		httputil.OK(w, groups)
	}
}

func (a *API) UpdateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
