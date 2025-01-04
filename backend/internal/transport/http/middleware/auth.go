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
				httputil.Forbidden(w)
				return
			}

			rights, ok := claims["rights"].([]interface{}) // JWT claims are usually unmarshaled as []interface{}
			if !ok {
				httputil.Forbidden(w)
				return
			}

			for _, userRight := range rights {
				if auth.IsEqual(userRight.(string), right) {
					next(w, r)
					return
				}
			}
			httputil.Forbidden(w)
			return
		}
	}
}

type RequestWithID interface {
	GetID() string
}

type ResourceIDSource int

const (
	FromBody ResourceIDSource = iota
	FromParam
)

type AuthConfig struct {
	ResourceType   string
	ResourceIDFrom ResourceIDSource
	RequiredRight  auth.Right

	// Used when ResourceIDFrom is FromParam
	ParamName string
	// Used when ResourceIDFrom is FromBody
	Body any
	// Whether to allow unauthenticated creation
	AllowCreate bool
}

func RequireAuth(config AuthConfig) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var resourceID string
			var isCreate bool

			switch config.ResourceIDFrom {
			case FromBody:
				if config.Body != nil && r.Body != nil {
					bodyData := config.Body.(RequestWithID)
					if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
						httputil.BadRequest(w)
						return
					}

					resourceID = bodyData.GetID()
					isCreate = resourceID == ""

					bodyBytes, _ := json.Marshal(bodyData)
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			case FromParam:
				if config.ParamName != "" {
					vars := mux.Vars(r)
					resourceID = vars[config.ParamName]
				}
			}

			// Only authenticate ownership if an ID is provided, otherwise this request is for creating a new resource
			authToken := r.Header.Get("Authorization")
			if authToken == "" {
				if isCreate && config.AllowCreate {
					next(w, r)
					return
				}
				httputil.Unauthorized(w)
				return
			}

			claims, err := auth.ValidateToken(authToken)
			if err != nil {
				httputil.Unauthorized(w)
				return
			}

			if !auth.HasOwnership(claims, config.ResourceType, resourceID) {
				httputil.Forbidden(w)
				return
			}

			if config.RequiredRight != "" {
				if !auth.HasRight(claims, config.RequiredRight) {
					httputil.Forbidden(w)
					return
				}
			}

			next(w, r)
		}
	}
}
