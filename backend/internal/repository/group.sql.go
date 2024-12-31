// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: group.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createGroupWithOwner = `-- name: CreateGroupWithOwner :one
WITH player_info AS (
    SELECT Players.id AS player_id, Players.display_name
    FROM Players
    WHERE Players.id = $1 -- player_id
),
new_group AS (
    INSERT INTO Groups (
        community_id,
        owner,
        region,
        gamemode,
        open,
        vanguards,
        duelists,
        strategists,
        platforms,
        voice_chat,
        mic
    ) 
    SELECT 
        $2,                  -- community_id
        player_info.display_name,  -- owner
        $3,                -- region
        $4,                -- gamemode
        $5,                -- open
        $6,                -- vanguards
        $7,                -- duelists
        $8,                -- strategists
        $9,                -- platforms
        $10,                -- voice_chat
        $11                -- mic
    FROM player_info
    RETURNING id
)
INSERT INTO GroupMembers (group_id, player_id, leader)
SELECT 
    new_group.id,
    player_info.player_id,
    TRUE
FROM new_group, player_info
RETURNING group_id
`

type CreateGroupWithOwnerParams struct {
	ID          int32       `json:"id"`
	CommunityID int32       `json:"community_id"`
	Region      string      `json:"region"`
	Gamemode    string      `json:"gamemode"`
	Open        bool        `json:"open"`
	Vanguards   pgtype.Int4 `json:"vanguards"`
	Duelists    pgtype.Int4 `json:"duelists"`
	Strategists pgtype.Int4 `json:"strategists"`
	Platforms   []string    `json:"platforms"`
	VoiceChat   pgtype.Bool `json:"voice_chat"`
	Mic         pgtype.Bool `json:"mic"`
}

func (q *Queries) CreateGroupWithOwner(ctx context.Context, arg CreateGroupWithOwnerParams) (string, error) {
	row := q.db.QueryRow(ctx, createGroupWithOwner,
		arg.ID,
		arg.CommunityID,
		arg.Region,
		arg.Gamemode,
		arg.Open,
		arg.Vanguards,
		arg.Duelists,
		arg.Strategists,
		arg.Platforms,
		arg.VoiceChat,
		arg.Mic,
	)
	var group_id string
	err := row.Scan(&group_id)
	return group_id, err
}
