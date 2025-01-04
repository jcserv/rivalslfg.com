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
	byPlayerID   = "/{playerID}"

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

func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(groups,
		middleware.RequireOwnership("group", "", &UpsertGroup{})(
			a.UpsertGroup(),
		),
	)

	r.HandleFunc(groups, a.UpsertGroup()).Methods(http.MethodPost)
	r.HandleFunc(group, a.GetGroupByID()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.GetGroups()).Methods(http.MethodGet)
	r.HandleFunc(groups, a.DeleteGroup()).Methods(http.MethodDelete)

	r.HandleFunc(playaz, a.CreatePlayer()).Methods(http.MethodPost)
	r.HandleFunc(player, a.JoinGroup()).Methods(http.MethodPost)
	r.HandleFunc(player,
		middleware.RequireRight(auth.RightUpdateGroup)(
			middleware.RequireOwnership("group", "id", nil)(
				a.RemovePlayer(),
			),
		),
	).Methods((http.MethodDelete))
}
