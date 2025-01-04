package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gorilla/mux"
	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/test"
	"github.com/jcserv/rivalslfg/internal/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestIntegration_UpsertGroup(t *testing.T) {
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
	t.Run("Should allow unauthenticated users to create a group", func(t *testing.T) {
		mockGroupService.EXPECT().UpsertGroup(gomock.Any(), gomock.Any()).Return("AAAA", nil)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", test.GetBody(
			mockUpsertGroup.ToMap(),
		))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should allow authenticated group owners to update group", func(t *testing.T) {
		mockGroupService.EXPECT().UpsertGroup(gomock.Any(), gomock.Any()).Return("AAAA", nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.RightUpdateGroup)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups", test.GetBody(
			getMockUpsertGroupDTOWithOverrides(
				&map[string]any{
					"id": "AAAA",
				},
			),
		))
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should prevent non-group owner from updating group", func(t *testing.T) {
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAB",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/groups", test.GetBody(
			getMockUpsertGroupDTOWithOverrides(
				&map[string]any{
					"id": "AAAA",
				},
			),
		))
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestIntegration_RemovePlayer(t *testing.T) {
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
	t.Run("Should allow group owners to kick other players", func(t *testing.T) {
		mockGroupService.EXPECT().RemovePlayerFromGroup(gomock.Any(), gomock.Any()).Return(nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
			"groupId":  "AAAA",
		}, auth.RightLeaveGroup)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"groupId":       "AAAA",
				"playerName":    "imphungky",
				"requesterName": "xZestence",
			},
		))
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Should allow group member leaving", func(t *testing.T) {
		mockGroupService.EXPECT().RemovePlayerFromGroup(gomock.Any(), gomock.Any()).Return(nil)
		token, _ := auth.GenerateToken("1", map[string]string{
			"playerId": "1",
		}, auth.RightLeaveGroup)
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/groups/AAAA/players/1", test.GetBody(
			map[string]interface{}{
				"groupId":       "AAAA",
				"playerName":    "xZestence",
				"requesterName": "xZestence",
			},
		))
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

var mockUpsertGroup = UpsertGroup{
	ID:       "",
	Owner:    "imphungky",
	Region:   "na",
	Gamemode: "competitive",
	Players: []repository.PlayerInGroup{
		{
			Name:     "imphungky",
			Leader:   true,
			Platform: "xb",
			Roles:    []string{"vanguard"},
			Rank:     "d3",
			Characters: []string{
				"Doctor Strange",
			},
			VoiceChat: true,
			Mic:       true,
		},
	},
	Open: true,
	RoleQueue: &repository.RoleQueue{
		Vanguards:   2,
		Duelists:    2,
		Strategists: 2,
	},
	GroupSettings: &repository.GroupSettings{
		Platforms: []string{
			"pc", "xb", "ps",
		},
		VoiceChat: true,
		Mic:       true,
	},
}

func getMockUpsertGroupDTOWithOverrides(overrides *map[string]any) map[string]any {
	return upsertGroupToMap(mockUpsertGroup, overrides)
}

func upsertGroupToMap(u UpsertGroup, overrides *map[string]any) map[string]any {
	out := map[string]any{}
	_ = mapstructure.Decode(u, &out)
	if overrides == nil {
		return out
	}

	for k, v := range *overrides {
		out[k] = v
	}
	return out
}
