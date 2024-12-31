package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

func (a *API) CreatePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input CreatePlayer
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			httputil.BadRequest(w)
			return
		}

		if err := input.Validate(); err != nil {
			httputil.BadRequest(w)
			return
		}

		playerID, err := a.playerService.CreatePlayer(ctx, input.ToParams())
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
			return
		}

		httputil.OK(w, map[string]int32{
			"id": playerID,
		})
	}
}

func (a *API) ReadPlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		playerID := vars["id"]
		if playerID == "" {
			httputil.BadRequest(w)
			return
		}

		playerIDVal, err := strconv.ParseInt(playerID, 10, 32)
		if err != nil {
			httputil.BadRequest(w)
			return
		}

		player, err := a.playerService.FindPlayer(ctx, int32(playerIDVal), "")
		if err != nil {
			httputil.InternalServerError(ctx, w, err)
			log.Error(ctx, err.Error())
			return
		}

		if player == nil {
			httputil.NotFound(w)
			return
		}

		httputil.OK(w, player)
	}
}
