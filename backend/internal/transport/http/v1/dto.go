package v1

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jcserv/rivalslfg/internal/repository"
)

type CreateGroup struct {
	PlayerID      int                       `json:"playerId"`
	Region        string                    `json:"region"`
	Gamemode      string                    `json:"gamemode"`
	Open          bool                      `json:"open"`
	RoleQueue     *repository.RoleQueue     `json:"roleQueue"`
	GroupSettings *repository.GroupSettings `json:"groupSettings"`
}

func (c *CreateGroup) Validate() error {
	if c.PlayerID <= 0 {
		return fmt.Errorf("PlayerID is required")
	}
	if c.Region == "" {
		return fmt.Errorf("region is required")
	}
	if c.Gamemode == "" {
		return fmt.Errorf("gamemode is required")
	}
	return nil
}

func (c *CreateGroup) ToParams() repository.CreateGroupWithOwnerParams {
	return repository.CreateGroupWithOwnerParams{
		ID:          int32(c.PlayerID),
		CommunityID: 1,
		Region:      c.Region,
		Gamemode:    c.Gamemode,
		Open:        c.Open,
		Vanguards:   pgtype.Int4{Int32: int32(c.RoleQueue.Vanguards), Valid: true},
		Duelists:    pgtype.Int4{Int32: int32(c.RoleQueue.Duelists), Valid: true},
		Strategists: pgtype.Int4{Int32: int32(c.RoleQueue.Strategists), Valid: true},
		Platforms:   c.GroupSettings.Platforms,
		VoiceChat:   pgtype.Bool{Bool: c.GroupSettings.VoiceChat, Valid: true},
		Mic:         pgtype.Bool{Bool: c.GroupSettings.Mic, Valid: true},
	}
}

type UpdatePlayer struct {
	DisplayName   string                    `json:"displayName"`
	Region        string                    `json:"region"`
	Platform      string                    `json:"platform"`
	Gamemode      string                    `json:"gamemode"`
	Roles         []string                  `json:"roles"`
	Rank          string                    `json:"rank"`
	Characters    []string                  `json:"characters"`
	VoiceChat     bool                      `json:"voiceChat"`
	Mic           bool                      `json:"mic"`
	RoleQueue     *repository.RoleQueue     `json:"roleQueue"`
	GroupSettings *repository.GroupSettings `json:"groupSettings"`
}

func (c *UpdatePlayer) Validate() error {
	if c.DisplayName == "" {
		return fmt.Errorf("display_name is required")
	}
	if c.Region == "" {
		return fmt.Errorf("region is required")
	}
	if c.Gamemode == "" {
		return fmt.Errorf("gamemode is required")
	}
	if c.Rank == "" {
		return fmt.Errorf("rank is required")
	}
	return nil
}

func (c *UpdatePlayer) ToParams(id int32) repository.UpsertPlayerParams {
	var vanguards, duelists, strategists pgtype.Int4
	if c.RoleQueue != nil {
		vanguards = pgtype.Int4{Int32: int32(c.RoleQueue.Vanguards), Valid: true}
		duelists = pgtype.Int4{Int32: int32(c.RoleQueue.Duelists), Valid: true}
		strategists = pgtype.Int4{Int32: int32(c.RoleQueue.Strategists), Valid: true}
	}

	var platforms []string
	var gVoiceChat, gMic pgtype.Bool
	if c.GroupSettings != nil {
		platforms = c.GroupSettings.Platforms
		gVoiceChat = pgtype.Bool{Bool: c.GroupSettings.VoiceChat, Valid: true}
		gMic = pgtype.Bool{Bool: c.GroupSettings.Mic, Valid: true}
	}

	return repository.UpsertPlayerParams{
		ID:          id,
		Name:        strings.ToLower(c.DisplayName),
		DisplayName: c.DisplayName,
		Region:      c.Region,
		Platform:    c.Platform,
		Gamemode:    c.Gamemode,
		Roles:       c.Roles,
		Rank:        repository.Rankid(c.Rank),
		Characters:  c.Characters,
		VoiceChat:   c.VoiceChat,
		Mic:         c.Mic,
		Vanguards:   vanguards,
		Duelists:    duelists,
		Strategists: strategists,
		Platforms:   platforms,
		GVoiceChat:  gVoiceChat,
		GMic:        gMic,
	}
}
