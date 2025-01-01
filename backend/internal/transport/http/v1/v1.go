package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
)

const (
	APIV1URLPath = "/api/v1/"
	byId         = "/{id}"

	groups    = APIV1URLPath + "groups"
	group     = groups + byId
	joinGroup = group + "/join"

	players = APIV1URLPath + "players"
	player  = players + byId
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
	r.HandleFunc(groups, a.UpsertGroup()).Methods(http.MethodPost)
	r.HandleFunc(group, a.GetGroupByID()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.GetGroups()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.DeleteGroup()).Methods(http.MethodDelete)
	r.HandleFunc(joinGroup, a.JoinGroup()).Methods(http.MethodPost)
}
