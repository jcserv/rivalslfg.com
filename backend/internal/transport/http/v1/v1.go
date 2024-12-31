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
	createGroup  = APIV1URLPath + "groups"
	getGroup     = APIV1URLPath + "groups/{id}"
	getGroupByID = APIV1URLPath + "groups"
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
	r.HandleFunc(createGroup, a.CreateGroup()).Methods(http.MethodPost)
	r.HandleFunc(getGroup, a.GetGroup()).Methods(http.MethodGet)
	r.HandleFunc(getGroupByID, a.GetGroups()).Methods(http.MethodGet)
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

func (a *API) GetGroup() http.HandlerFunc {
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

		httputil.OK(w, group)
	}
}

// TODO: Add pagination
func (a *API) GetGroups() http.HandlerFunc {
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
