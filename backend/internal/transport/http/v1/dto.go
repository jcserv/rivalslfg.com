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
