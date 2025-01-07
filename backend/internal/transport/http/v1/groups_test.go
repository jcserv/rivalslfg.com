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

func TestIntegration_CreateGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := mux.NewRouter()
	mockGroupService := mocks.NewMockIGroup(ctrl)
	mockPlayerService := mocks.NewMockIPlayer(ctrl)

	a := NewAPI(
		&Dependencies{
			GroupService:  mockGroupService,
			PlayerService: mockPlayerService,
		},
	)
	a.RegisterRoutes(r)
	t.Run("Should allow unauthenticated users to create a group", func(t *testing.T) {
		mockGroupService.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Return(repository.CreateGroupRow{
			GroupID:  "AAAA",
			PlayerID: 1,
		}, nil)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", test.GetBody(
			map[string]interface{}{
				"owner":    "imphungky",
				"passcode": "abcd",
				"gamemode": "competitive",
				"region":   "na",
				"platform": "xb",
				"role":     "vanguard",
				"rankId":   "d3",
				"characters": []string{
					"Doctor Strange",
				},
				"voiceChat":   true,
				"mic":         true,
				"vanguards":   2,
				"duelists":    2,
				"strategists": 2,
			},
		))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestIntegration_GetGroupByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := mux.NewRouter()
	mockGroupService := mocks.NewMockIGroup(ctrl)
	mockPlayerService := mocks.NewMockIPlayer(ctrl)

	a := NewAPI(
		&Dependencies{
			GroupService:  mockGroupService,
			PlayerService: mockPlayerService,
		},
	)
	a.RegisterRoutes(r)
	t.Parallel()
	t.Run("Should return group if exists and is open", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA", false).Return(&repository.GroupWithPlayers{
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
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA", false).Return(nil, fmt.Errorf("unexpected err"))
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
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA", false).Return(nil, nil)
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
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA", false).Return(&repository.GroupWithPlayers{
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
		mockGroupService.EXPECT().GetGroupByID(gomock.Any(), "AAAA", false).Return(&repository.GroupWithPlayers{
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

func TestIntegration_GetGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := mux.NewRouter()
	mockGroupService := mocks.NewMockIGroup(ctrl)
	mockPlayerService := mocks.NewMockIPlayer(ctrl)

	a := NewAPI(
		&Dependencies{
			GroupService:  mockGroupService,
			PlayerService: mockPlayerService,
		},
	)
	a.RegisterRoutes(r)
	t.Parallel()
	t.Run("Should return 200 for simple query", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroups(gomock.Any(), gomock.Any()).Return([]repository.GroupWithPlayers{}, int32(0), nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("Should return 200 on valid, full query", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroups(gomock.Any(), gomock.Any()).Return([]repository.GroupWithPlayers{
			{
				GroupDTO: repository.GroupDTO{
					ID:            "AAAA",
					CommunityID:   int32(1),
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
			},
		}, int32(100), nil)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups", nil)
		q := req.URL.Query()
		q.Add("filters", "region eq \"na\" and gamemode eq \"competitive\" and open eq true")
		q.Add("sort", "size")
		q.Add("limit", "1")
		q.Add("offset", "0")
		q.Add("count", "true")
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "100", rec.Header().Get("X-Total-Count"))
	})
	t.Run("Should return 400 on invalid query param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups", nil)
		q := req.URL.Query()
		q.Add("limit", "-1")
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("Should return 500 on unexpected error", func(t *testing.T) {
		mockGroupService.EXPECT().GetGroups(gomock.Any(), gomock.Any()).Return(nil, int32(0), fmt.Errorf("unexpected error"))
		req := httptest.NewRequest(http.MethodGet, "/api/v1/groups", nil)
		q := req.URL.Query()
		q.Add("filters", "region eq \"na\"")
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
