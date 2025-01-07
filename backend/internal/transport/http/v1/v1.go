package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
)

const (
	APIV1URLPath = "/api/v1/"
	byId         = "/{id}"
	byPlayerID   = "/{playerId}"

	groups       = APIV1URLPath + "groups"
	group        = groups + byId
	groupDetails = group + "/details"

	players = APIV1URLPath + "players"

	groupMembers = group + "/players"
	groupMember  = groupMembers + byPlayerID
)

type API struct {
	groupService  services.IGroup
	playerService services.IPlayer
}

type Dependencies struct {
	GroupService  services.IGroup
	PlayerService services.IPlayer
}

func NewAPI(deps *Dependencies) *API {
	return &API{
		groupService:  deps.GroupService,
		playerService: deps.PlayerService,
	}
}

// RegisterRoutes registers the routes for the V1 API.
func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(groups, a.CreateGroup()).Methods(http.MethodPost)
	r.HandleFunc(group, a.GetGroupByID()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.GetGroups()).Methods(http.MethodGet)
	r.HandleFunc(groupMember, a.JoinGroup()).Methods(http.MethodPost)

	r.HandleFunc(players, a.CreatePlayer()).Methods(http.MethodPost)

	// r.HandleFunc(groups,
	// 	middleware.RequireRight(auth.RightDeleteGroup)(
	// 		a.DeleteGroup(),
	// 	),
	// ).Methods(http.MethodDelete)
	// r.HandleFunc(player,
	// 	middleware.RequireRight(
	// 		auth.RightLeaveGroup,
	// 	)(a.RemovePlayer()),
	// ).Methods(http.MethodDelete)
}
