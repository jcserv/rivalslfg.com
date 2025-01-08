package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/test"
	"github.com/jcserv/rivalslfg/internal/test/mocks"
	"github.com/jcserv/rivalslfg/internal/transport/http/reqCtx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestIntegration_JoinGroup(t *testing.T) {
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
	t.Run("Should allow unauthenticated users to join a group", func(t *testing.T) {
		mockPlayerService.EXPECT().JoinGroup(gomock.Any(), gomock.Any()).Return(int32(1), nil)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"name":     "imphungky",
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
	t.Run("Should return 400 if required field is missing/empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"name":     "imphungky",
				"passcode": "abcd",
				"gamemode": "",
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
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("Should return 400 if serviceErr.BadRequest returned", func(t *testing.T) {
		mockPlayerService.EXPECT().JoinGroup(gomock.Any(), gomock.Any()).Return(int32(0), services.NewError(http.StatusBadRequest, "bad request", nil))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"name":     "imphungky",
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
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("Should return 403 if serviceErr.Forbidden returned", func(t *testing.T) {
		mockPlayerService.EXPECT().JoinGroup(gomock.Any(), gomock.Any()).Return(int32(0), services.NewError(http.StatusForbidden, "forbidden", nil))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"name":     "imphungky",
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
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
	t.Run("Should return 404 if serviceErr.NotFound returned", func(t *testing.T) {
		mockPlayerService.EXPECT().JoinGroup(gomock.Any(), gomock.Any()).Return(int32(0), services.NewError(http.StatusNotFound, "not found", nil))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"name":     "imphungky",
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
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
	t.Run("Should return 500 if unexpected error", func(t *testing.T) {
		mockPlayerService.EXPECT().JoinGroup(gomock.Any(), gomock.Any()).Return(int32(0), fmt.Errorf("unexpected error"))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"name":     "imphungky",
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
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}

func TestIntegration_RemovePlayer(t *testing.T) {
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
	t.Run("Should allow group owner to remove regular member", func(t *testing.T) {
		mockPlayerService.EXPECT().RemovePlayer(gomock.Any(), gomock.Any()).Return("200", nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/2", nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.GroupOwnerRights...)

		req.Header.Set("Authorization", token)
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			PlayerID: 1,
			GroupID:  "AAAA",
			Token:    token,
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should allow group member to remove themselves", func(t *testing.T) {
		mockPlayerService.EXPECT().RemovePlayer(gomock.Any(), gomock.Any()).Return("200", nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/2", nil)
		token, _ := auth.GenerateToken("2", map[string]string{
			"playerId": "2",
			"groupId":  "AAAA",
		}, auth.GroupMemberRights...)

		req.Header.Set("Authorization", token)
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			PlayerID: 2,
			GroupID:  "AAAA",
			Token:    token,
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should return 204 when removing last member", func(t *testing.T) {
		mockPlayerService.EXPECT().RemovePlayer(gomock.Any(), gomock.Any()).Return("204", nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/1", nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.GroupOwnerRights...)

		req.Header.Set("Authorization", token)
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			PlayerID: 1,
			GroupID:  "AAAA",
			Token:    token,
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Should return 404 if player not found in group", func(t *testing.T) {
		mockPlayerService.EXPECT().RemovePlayer(gomock.Any(), gomock.Any()).Return("", services.NewError(http.StatusNotFound, "Player not found.", nil))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/999", nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.GroupOwnerRights...)

		req.Header.Set("Authorization", token)
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			PlayerID: 1,
			GroupID:  "AAAA",
			Token:    token,
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Should return 403 if non-owner attempts to remove another player", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/3", nil)
		token, _ := auth.GenerateToken("2", map[string]string{
			"playerId": "2",
			"groupId":  "AAAA",
		}, auth.GroupMemberRights...)

		req.Header.Set("Authorization", token)
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			PlayerID: 2,
			GroupID:  "AAAA",
			Token:    token,
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Should return 403 if unauthorized user attempts to remove player", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/1", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Should return 500 if unexpected error occurs", func(t *testing.T) {
		mockPlayerService.EXPECT().RemovePlayer(gomock.Any(), gomock.Any()).Return("", fmt.Errorf("unexpected error"))

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/2", nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.GroupOwnerRights...)

		req.Header.Set("Authorization", token)
		req = reqCtx.WithAuthInfo(req, &reqCtx.AuthInfo{
			PlayerID: 1,
			GroupID:  "AAAA",
			Token:    token,
		})
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
