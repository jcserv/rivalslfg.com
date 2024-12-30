package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/transport/rest/httputil"
)

const (
	APIV1URLPath = "/api/v1/"
	GetGroup     = APIV1URLPath + "group/:id"
	GetGroups    = APIV1URLPath + "groups"
)

type API struct {
	conn *pgx.Conn
}

func NewAPI(conn *pgx.Conn) *API {
	return &API{
		conn: conn,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(GetGroup, a.GetGroup()).Methods(http.MethodGet)
	r.HandleFunc(GetGroups, a.GetGroups()).Methods(http.MethodGet)
}

func (a *API) GetGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.OK(w, map[string]string{
			"message": "get group",
		})
	}
}

func (a *API) GetGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		repo := repository.New(a.conn)
		players, _ := repo.FindAllPlayers(ctx)
		httputil.OK(w, map[string][]repository.Player{
			"players": players,
		})
	}
}
