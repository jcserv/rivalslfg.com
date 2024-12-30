package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/transport/rest/httputil"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

const (
	APIV1URLPath = "/api/v1/"
	GetGroup     = APIV1URLPath + "groups/{id}"
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
		ctx := r.Context()
		vars := mux.Vars(r)
		groupID := vars["id"]
		if groupID == "" {
			httputil.BadRequest(w)
			return
		}

		repo := repository.New(a.conn)

		group, err := repo.GetGroupByID(ctx, groupID)
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
			return
		}

		httputil.OK(w, group)
	}
}

func (a *API) GetGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		repo := repository.New(a.conn)

		groups, err := repo.FindAllGroups(ctx)
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
			return
		}

		httputil.OK(w, groups)
	}
}
