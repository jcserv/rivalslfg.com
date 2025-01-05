package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/transport/http/middleware"
)

const (
	APIV1URLPath = "/api/v1/"
	byId         = "/{id}"
	byPlayerID   = "/{playerId}"

	groups = APIV1URLPath + "groups"
	group  = groups + byId

	// todo
	playaz = APIV1URLPath + "players"

	players = group + "/players"
	player  = players + byPlayerID
)

type API struct {
	authService  services.IAuth
	groupService services.IGroup
}

func NewAPI(authService services.IAuth, groupService services.IGroup) *API {
	return &API{
		authService:  authService,
		groupService: groupService,
	}
}

// RegisterRoutes registers the routes for the V1 API.
func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(group, a.GetGroupByID()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.GetGroups()).Methods(http.MethodGet)
	r.HandleFunc(player, a.JoinGroup()).Methods(http.MethodPost)
	r.HandleFunc(playaz, a.CreatePlayer()).Methods(http.MethodPost)

	r.HandleFunc(groups,
		middleware.RequireRight(auth.RightDeleteGroup)(
			a.DeleteGroup(),
		),
	).Methods(http.MethodDelete)
	r.HandleFunc(player,
		middleware.RequireRight(
			auth.RightLeaveGroup,
		)(a.RemovePlayer()),
	).Methods(http.MethodDelete)

	r.HandleFunc(groups,
		middleware.RequireAuth(middleware.AuthConfig{
			ResourceType:   "group",
			ResourceIDFrom: middleware.FromBody,
			RequiredRight:  auth.RightUpdateGroup,
			Body:           &UpsertGroup{},
			AllowCreate:    true,
		})(a.UpsertGroup()),
	)
}
