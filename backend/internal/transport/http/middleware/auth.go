package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
)

func RequireRight(right auth.Right) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if authToken == "" {
				httputil.Unauthorized(w)
				return
			}

			claims, err := auth.ValidateToken(authToken)
			if err != nil {
				httputil.Unauthorized(w)
				return
			}

			rights, ok := claims["rights"].([]interface{}) // JWT claims are usually unmarshaled as []interface{}
			if !ok {
				httputil.Unauthorized(w)
				return
			}

			for _, userRight := range rights {
				if auth.IsEqual(userRight.(string), right) {
					next(w, r)
					return
				}
			}
			httputil.Unauthorized(w)
			return
		}
	}
}

type RequestWithID interface {
	GetID() string
}

func RequireOwnership(resourceType string, paramName string, body any) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var resourceID string

			if body != nil && r.Body != nil {
				bodyData := body.(RequestWithID)
				if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
					httputil.BadRequest(w)
					return
				}

				// Only authenticate ownership if an ID is provided, otherwise this request is for creating a new resource
				if bodyData.GetID() != "" {
					resourceID = bodyData.GetID()
				}

				// IMPORTANT: restore body for later use
				bodyBytes, _ := json.Marshal(bodyData)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			authToken := r.Header.Get("Authorization")
			// Allow unauthenticated users to create a group
			if authToken == "" && resourceID == "" {
				next(w, r)
				return
			}

			claims, err := auth.ValidateToken(authToken)
			if err != nil {
				httputil.Unauthorized(w)
				return
			}

			if paramName != "" {
				vars := mux.Vars(r)
				resourceID = vars[paramName]
			}

			if !auth.HasOwnership(claims, resourceType, resourceID) {
				httputil.Forbidden(w)
				return
			}

			next(w, r)
		}
	}
}
