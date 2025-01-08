package v1

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroup_Validate(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
			Vanguards:   2,
			Duelists:    2,
			Strategists: 2,
			VoiceChat:   true,
			Mic:         true,
			Open:        true,
			Platforms:   []string{"pc", "ps"},
		}
		err := input.validate()
		assert.NoError(t, err)
	})

	t.Run("Should validate owner", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "owner is required")
	})

	t.Run("Should validate platform", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "invalid",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "platform invalid is not supported")
	})

	t.Run("Should validate role", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "invalid",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role invalid is not supported")
	})

	t.Run("Should validate rank", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "invalid",
			Characters: []string{
				"Doctor Strange",
			},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid rank invalid")
	})

	t.Run("Should validate region", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "invalid",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "region invalid is not supported")
	})

	t.Run("Should validate gamemode", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "invalid",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gamemode invalid is not supported")
	})

	t.Run("Should validate role queue", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
			Vanguards:   -1,
			Duelists:    7,
			Strategists: 7,
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be between 0 and 6")
	})

	t.Run("Should validate platforms", func(t *testing.T) {
		input := CreateGroup{
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
			Platforms: []string{"invalid"},
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "one or more provided platforms")
	})
}

func TestCreateGroup_Parse(t *testing.T) {
	t.Run("Should parse input to repository params", func(t *testing.T) {
		input := CreateGroup{
			PlayerID: 1,
			GroupID:  "AAAA",
			Owner:    "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "Vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
			VoiceChat: true,
			Mic:       true,
		}

		result, err := input.Parse()
		assert.NoError(t, err)
		assert.Equal(t, int32(1), result.PlayerID)
		assert.Equal(t, "AAAA", result.GroupID)
		assert.Equal(t, "imphungky", result.Owner)
		assert.Equal(t, "na", result.Region)
		assert.Equal(t, "competitive", result.Gamemode)
		assert.Equal(t, "vanguard", result.Role)
		assert.Equal(t, "pc", result.Platform)
		assert.Equal(t, int32(40), result.RankVal) // d3 = 40
		assert.Equal(t, []string{"Doctor Strange"}, result.Characters)
		assert.True(t, result.VoiceChat)
		assert.True(t, result.Mic)
		assert.Equal(t, pgtype.Bool{Bool: false, Valid: true}, result.GroupMic)
		assert.Equal(t, pgtype.Bool{Bool: false, Valid: true}, result.GroupVoiceChat)
	})
}

func TestJoinGroup_Validate(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
			VoiceChat: true,
			Mic:       true,
		}
		err := input.validate()
		assert.NoError(t, err)
	})

	t.Run("Should validate groupId", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "",
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "groupId is required")
	})

	t.Run("Should validate name", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "playerName is required")
	})

	t.Run("Should validate gamemode", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "invalid",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gamemode invalid is not supported")
	})

	t.Run("Should validate region", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "imphungky",
			Region:   "invalid",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "d3",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "region invalid is not supported")
	})

	t.Run("Should validate platform", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "invalid",
			RankID:   "d3",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "platform invalid is not supported")
	})

	t.Run("Should validate role", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "invalid",
			Platform: "pc",
			RankID:   "d3",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role invalid is not supported")
	})

	t.Run("Should validate rankId", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "vanguard",
			Platform: "pc",
			RankID:   "invalid",
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rankId invalid is invalid")
	})
}

func TestJoinGroup_Parse(t *testing.T) {
	t.Run("Should parse input to repository params", func(t *testing.T) {
		input := JoinGroup{
			GroupID:  "AAAA",
			PlayerID: 1,
			Name:     "imphungky",
			Region:   "na",
			Gamemode: "competitive",
			Role:     "Vanguard",
			Platform: "pc",
			RankID:   "d3",
			Characters: []string{
				"Doctor Strange",
			},
			VoiceChat: true,
			Mic:       true,
		}

		result, err := input.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "AAAA", result.GroupID)
		assert.Equal(t, int32(1), result.PlayerID)
		assert.Equal(t, "imphungky", result.Name)
		assert.Equal(t, "na", result.Region)
		assert.Equal(t, "competitive", result.Gamemode)
		assert.Equal(t, "vanguard", result.Role)
		assert.Equal(t, "pc", result.Platform)
		assert.Equal(t, int32(40), result.RankVal) // d3 = 40
		assert.Equal(t, []string{"Doctor Strange"}, result.Characters)
		assert.True(t, result.VoiceChat)
		assert.True(t, result.Mic)
	})
}

func TestRemovePlayer_Validate(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := RemovePlayer{
			GroupID:          "AAAA",
			PlayerToRemoveID: 1,
		}
		err := input.validate()
		assert.NoError(t, err)
	})

	t.Run("Should validate groupId", func(t *testing.T) {
		input := RemovePlayer{
			GroupID:          "",
			PlayerToRemoveID: 1,
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "groupId is required")
	})

	t.Run("Should validate playerId", func(t *testing.T) {
		input := RemovePlayer{
			GroupID:          "AAAA",
			PlayerToRemoveID: 0,
		}
		err := input.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "playerId is required")
	})
}

func TestRemovePlayer_Parse(t *testing.T) {
	t.Run("Should parse input to repository params", func(t *testing.T) {
		input := RemovePlayer{
			GroupID:          "AAAA",
			PlayerToRemoveID: 1,
		}

		result, err := input.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "AAAA", result.GroupID)
		assert.Equal(t, int32(1), result.PlayerID)
	})

	t.Run("Should handle normal int to int32 conversion", func(t *testing.T) {
		input := RemovePlayer{
			GroupID:          "AAAA",
			PlayerToRemoveID: 2147483647, // max int32
		}

		result, err := input.Parse()
		assert.NoError(t, err)
		assert.Equal(t, "AAAA", result.GroupID)
		assert.Equal(t, int32(2147483647), result.PlayerID)
	})
}
