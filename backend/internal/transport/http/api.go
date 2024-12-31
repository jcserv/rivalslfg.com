package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	v1 "github.com/jcserv/rivalslfg/internal/transport/http/v1"
)

const (
	HealthCheck = "/health/system"
)

type API struct {
	V1API *v1.API
}

func NewAPI(groupService *services.GroupService, playerService *services.PlayerService) *API {
	return &API{
		V1API: v1.NewAPI(groupService, playerService),
	}
}

func (a *API) RegisterRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(LogIncomingRequests())
	a.V1API.RegisterRoutes(r)
	r.HandleFunc(HealthCheck, a.HealthCheck()).Methods(http.MethodGet)
	return r
}

func (a *API) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.OK(w, nil)
	}
}
