package v1

import (
	"fmt"

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

type CreatePlayer struct {
	// Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Region      string      `json:"region"`
	Platform    string      `json:"platform"`
	Gamemode    string      `json:"gamemode"`
	Roles       []string    `json:"roles"`
	Rank        string      `json:"rank"`
	Characters  []string    `json:"characters"`
	VoiceChat   bool        `json:"voice_chat"`
	Mic         bool        `json:"mic"`
	Vanguards   pgtype.Int4 `json:"vanguards"`
	Duelists    pgtype.Int4 `json:"duelists"`
	Strategists pgtype.Int4 `json:"strategists"`
	Platforms   []string    `json:"platforms"`
	GVoiceChat  pgtype.Bool `json:"g_voice_chat"`
	GMic        pgtype.Bool `json:"g_mic"`
}

func (c *CreatePlayer) Validate() error {
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

func (c *CreatePlayer) ToParams() repository.CreatePlayerParams {
	return repository.CreatePlayerParams{
		Name:        c.DisplayName,
		DisplayName: c.DisplayName,
		Region:      c.Region,
		Platform:    c.Platform,
		Gamemode:    c.Gamemode,
		Roles:       c.Roles,
		Rank:        repository.Rankid(c.Rank),
		Characters:  c.Characters,
		VoiceChat:   c.VoiceChat,
		Mic:         c.Mic,
		Vanguards:   c.Vanguards,
		Duelists:    c.Duelists,
		Strategists: c.Strategists,
		Platforms:   c.Platforms,
		GVoiceChat:  c.GVoiceChat,
		GMic:        c.GMic,
	}
}
