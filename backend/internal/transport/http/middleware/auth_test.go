package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequireRight(t *testing.T) {
	t.Run("HasRight", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":    "user123",
			"rights": []string{"user:read"},
		}
		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "claims", claims))
		rec := httptest.NewRecorder()

		handler := RequireRight("user:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("NoRight", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":    "user123",
			"rights": []string{"user:write"},
		}
		req := httptest.NewRequest("GET", "/", nil)
		req = req.WithContext(context.WithValue(req.Context(), "claims", claims))
		rec := httptest.NewRecorder()

		handler := RequireRight("user:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
