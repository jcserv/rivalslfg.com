package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/services"
	"github.com/jcserv/rivalslfg/internal/test"
	"github.com/jcserv/rivalslfg/internal/test/mocks"
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
