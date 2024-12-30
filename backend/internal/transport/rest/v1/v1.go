package v1

import (
	"github.com/gorilla/mux"
)

const (
	APIV1URLPath = "/api/v1/"
)

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a *API) RegisterRoutes(r *mux.Router) {
}
