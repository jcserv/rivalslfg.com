package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
)

func InitRequestContext() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if authToken == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := auth.ValidateToken(authToken)
			if err != nil {
				httputil.Forbidden(w)
				return
			}
			next.ServeHTTP(w, reqCtx.Init(r, claims, authToken))
		})
	}
}

func RequireRight(right auth.Right) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authToken := reqCtx.GetToken(ctx)
			if authToken == "" {
				httputil.Unauthorized(w)
				return
			}

			claims, err := auth.ValidateToken(authToken)
			if err != nil {
				httputil.Forbidden(w)
				return
			}

			rights, ok := claims["rights"].([]interface{})
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
						httputil.BadRequest(w, err)
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
			ctx := r.Context()
			authToken := reqCtx.GetToken(ctx)
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

			if isCreate && config.AllowCreate && auth.HasNotCreatedGroup(claims) {
				next(w, reqCtx.Init(r, claims, authToken))
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
