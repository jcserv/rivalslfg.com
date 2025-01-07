package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/test"
	"github.com/jcserv/rivalslfg/internal/test/mocks"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestIntegration_GetGroupByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := mux.NewRouter()
	mockAuthService := mocks.NewMockIAuth(ctrl)
	mockGroupService := mocks.NewMockIGroup(ctrl)

	a := NewAPI(
		mockAuthService,
		mockGroupService,
	)
	a.RegisterRoutes(r)
	t.Parallel()
	t.Run("Should return group if exists and is open", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA").Return(&repository.GroupWithPlayers{
			GroupDTO: repository.GroupDTO{
				ID:            "AAAA",
				CommunityID:   1,
				Owner:         "imphungky",
				Region:        "na",
				Gamemode:      "competitive",
				Open:          true,
				Passcode:      "",
				RoleQueue:     nil,
				GroupSettings: nil,
				LastActiveAt:  time.Time{},
			},
			Name:    "imphungky's Group",
			Size:    0,
			Players: []repository.PlayerInGroup{},
		}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups/AAAA", test.GetBody(
			map[string]interface{}{
				"id": "AAAA",
			},
		))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("Should return 500 if unexpected error", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA").Return(nil, fmt.Errorf("unexpected err"))
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups/AAAA", test.GetBody(
			map[string]interface{}{
				"id": "AAAA",
			},
		))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("Should return 404 if null is returned from service", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA").Return(nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups/AAAA", test.GetBody(
			map[string]interface{}{
				"id": "AAAA",
			},
		))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
	t.Run("Should return 200 if group is not open and user is authorized to view", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA").Return(&repository.GroupWithPlayers{
			GroupDTO: repository.GroupDTO{
				ID:            "AAAA",
				CommunityID:   1,
				Owner:         "imphungky",
				Region:        "na",
				Gamemode:      "competitive",
				Open:          false,
				Passcode:      "",
				RoleQueue:     nil,
				GroupSettings: nil,
				LastActiveAt:  time.Time{},
			},
			Name:    "imphungky's Group",
			Size:    0,
			Players: []repository.PlayerInGroup{},
		}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups/AAAA", test.GetBody(
			map[string]interface{}{
				"id": "AAAA",
			},
		))
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			GroupID: "AAAA",
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("Should return 403 if group is not open and user is not authorized to view", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA").Return(&repository.GroupWithPlayers{
			GroupDTO: repository.GroupDTO{
				ID:            "AAAA",
				CommunityID:   1,
				Owner:         "imphungky",
				Region:        "na",
				Gamemode:      "competitive",
				Open:          false,
				Passcode:      "",
				RoleQueue:     nil,
				GroupSettings: nil,
				LastActiveAt:  time.Time{},
			},
			Name:    "imphungky's Group",
			Size:    0,
			Players: []repository.PlayerInGroup{},
		}, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups/AAAA", test.GetBody(
			map[string]interface{}{
				"id": "AAAA",
			},
		))
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			GroupID: "AAAB",
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
