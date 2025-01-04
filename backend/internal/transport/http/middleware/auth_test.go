package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestRequireRight(t *testing.T) {
	t.Run("Should allow if requested right is provided", func(t *testing.T) {
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.RightReadUser)
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		handler := RequireRight(auth.RightReadUser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("Should not allow if requested right is not provided", func(t *testing.T) {
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
		})
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		handler := RequireRight(auth.RightDeleteGroup)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
	t.Run("Should not allow if user does not have token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		handler := RequireRight(auth.RightReadUser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
