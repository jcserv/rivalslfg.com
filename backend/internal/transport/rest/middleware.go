package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

func LogIncomingRequests() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log.Info(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		})
	}
}
