package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TODO: Compute requirementsMet, add total count
const getGroups = `
SELECT 
    g.id,
    community_id,
    owner,
    g.region,
    g.gamemode,
    open,
    jsonb_build_object(
        'vanguards', g.vanguards,
        'duelists', g.duelists,
        'strategists', g.strategists
    ) AS role_queue,
    jsonb_build_object(
        'platforms', g.platforms,
        'voiceChat', g.voice_chat,
        'mic', g.mic
    ) AS group_settings,
    players,
    jsonb_array_length(players) AS size,
	last_active_at
FROM Groups g
WHERE ($1 = '' OR g.region = $1)
  AND ($2 = '' OR g.gamemode = $2)
  AND ($3 = '' OR
    CASE 
        WHEN LOWER($3) = 'true' THEN g.open = true
        WHEN LOWER($3) = 'false' THEN g.open = false
        ELSE TRUE -- Ignore invalid values and don't filter
    END
  )
ORDER BY 
    CASE WHEN $4 = 'asc' THEN jsonb_array_length(players) END ASC,
    CASE WHEN $4 = 'desc' THEN jsonb_array_length(players) END DESC
LIMIT $5 OFFSET $6;
`

type GetGroupsParams struct {
	RegionFilter   string `json:"regionFilter"`
	GamemodeFilter string `json:"gamemodeFilter"`
	OpenFilter     string `json:"openFilter"` // openFilter is a string to account for when we don't want to filter
	SizeSort       string `json:"sizeSort"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
}

func (q *Queries) GetGroups(ctx context.Context, arg GetGroupsParams) ([]GroupWithPlayers, error) {
	rows, err := q.db.Query(ctx, getGroups, arg.RegionFilter, arg.GamemodeFilter, arg.OpenFilter, arg.SizeSort, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]GroupWithPlayers, 0)
	for rows.Next() {
		var g GroupWithPlayers
		if err := rows.Scan(
			&g.ID,
			&g.CommunityID,
			&g.Owner,
			&g.Region,
			&g.Gamemode,
			&g.Open,
			&g.RoleQueue,
			&g.GroupSettings,
			&g.Players,
			&g.Size,
			&g.LastActiveAt,
		); err != nil {
			return nil, err
		}
		g.Name = fmt.Sprintf("%s's Group", g.Owner)
		result = append(result, g)
	}
	return result, nil
}

const getGroupByID = `
SELECT 
	g.id,
	community_id,
	owner,
	g.region,
	g.gamemode,
	open,
	jsonb_build_object(
		'vanguards', g.vanguards,
		'duelists', g.duelists,
		'strategists', g.strategists
	) AS role_queue,
	 jsonb_build_object(
		'platforms', g.platforms,
		'voiceChat', g.voice_chat,
		'mic', g.mic
	) AS group_settings,
	players,
	jsonb_array_length(players, 1) AS size,
	last_active_at
FROM Groups g
WHERE g.id = $1`

func (q *Queries) GetGroupByID(ctx context.Context, id string) (*GroupWithPlayers, error) {
	var g GroupWithPlayers
	err := q.db.QueryRow(ctx, getGroupByID, id).Scan(
		&g.ID,
		&g.CommunityID,
		&g.Owner,
		&g.Region,
		&g.Gamemode,
		&g.Open,
		&g.RoleQueue,
		&g.GroupSettings,
		&g.Players,
		&g.Size,
		&g.LastActiveAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	g.Name = fmt.Sprintf("%s's Group", g.Owner)
	return &g, nil
}
