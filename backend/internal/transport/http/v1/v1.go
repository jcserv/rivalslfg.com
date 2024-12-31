package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
)

const (
	APIV1URLPath = "/api/v1/"
	groups       = APIV1URLPath + "groups"
	group        = APIV1URLPath + "groups/{id}"
	joinGroup    = group + "/join"

	player = APIV1URLPath + "players/{id}"
)

type API struct {
	groupService  *services.GroupService
	playerService *services.PlayerService
}

func NewAPI(groupService *services.GroupService, playerService *services.PlayerService) *API {
	return &API{
		groupService:  groupService,
		playerService: playerService,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(groups, a.CreateGroup()).Methods(http.MethodPost)
	r.HandleFunc(group, a.ReadGroup()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.ReadGroups()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.UpdateGroup()).Methods(http.MethodPut)
	r.HandleFunc(groups, a.DeleteGroup()).Methods(http.MethodDelete)
	r.HandleFunc(joinGroup, a.JoinGroup()).Methods(http.MethodPost)

	r.HandleFunc(player, a.ReadPlayer()).Methods(http.MethodGet)
}
