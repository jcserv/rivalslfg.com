package v1

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jcserv/rivalslfg/internal/repository"
)

type UpsertGroup struct {
	ID            string                     `json:"id"`
	Owner         string                     `json:"owner"`
	Region        string                     `json:"region"`
	Gamemode      string                     `json:"gamemode"`
	Players       []repository.PlayerInGroup `json:"players"`
	Open          bool                       `json:"open"`
	RoleQueue     *repository.RoleQueue
	GroupSettings *repository.GroupSettings
}

func (c *UpsertGroup) Parse() (*repository.UpsertGroupParams, error) {
	if c.Owner == "" {
		return nil, fmt.Errorf("owner is required")
	}
	if c.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	if c.Gamemode == "" {
		return nil, fmt.Errorf("gamemode is required")
	}
	if len(c.Players) == 0 {
		return nil, fmt.Errorf("one or more players are required")
	}

	params := &repository.UpsertGroupParams{
		ID:       c.ID,
		Owner:    c.Owner,
		Region:   c.Region,
		Gamemode: c.Gamemode,
		Players:  c.Players,
		Open:     c.Open,
	}

	if c.RoleQueue != nil {
		params.Vanguards = pgtype.Int4{Int32: int32(c.RoleQueue.Vanguards), Valid: true}
		params.Duelists = pgtype.Int4{Int32: int32(c.RoleQueue.Duelists), Valid: true}
		params.Strategists = pgtype.Int4{Int32: int32(c.RoleQueue.Strategists), Valid: true}
	}
	if c.GroupSettings != nil {
		params.Platforms = c.GroupSettings.Platforms
		params.VoiceChat = pgtype.Bool{Bool: c.GroupSettings.VoiceChat, Valid: true}
		params.Mic = pgtype.Bool{Bool: c.GroupSettings.Mic, Valid: true}
	}

	return params, nil
}
