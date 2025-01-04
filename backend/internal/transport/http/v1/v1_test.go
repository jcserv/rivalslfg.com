package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthenticatedRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := mux.NewRouter()
	tokenService := auth.NewTokenService("test-secret")
	a := NewAPI(
		mocks.NewMockIAuth(ctrl),
		mocks.NewMockIGroup(ctrl),
	)

	a.RegisterRoutes(r)

	t.Run("Group owner can delete their group", func(t *testing.T) {
		token, _ := tokenService.GenerateToken("user123", true)
		req := httptest.NewRequest("DELETE", "/api/v1/groups/AAAA", nil)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
