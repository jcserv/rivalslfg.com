package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const findAllGroups = `
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
    players
FROM Groups g`

func (q *Queries) FindAllGroups(ctx context.Context) ([]GroupWithPlayers, error) {
	rows, err := q.db.Query(ctx, findAllGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []GroupWithPlayers
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
		); err != nil {
			return nil, err
		}
		g.Name = fmt.Sprintf("%s's Group", g.Owner)
		result = append(result, g)
	}
	return result, nil
}

const findGroupByID = `
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
	last_active_at,
FROM Groups g
WHERE g.id = $1`

func (q *Queries) FindGroupByID(ctx context.Context, id string) (*GroupWithPlayers, error) {
	var g GroupWithPlayers
	err := q.db.QueryRow(ctx, findGroupByID, id).Scan(
		&g.ID,
		&g.CommunityID,
		&g.Owner,
		&g.Region,
		&g.Gamemode,
		&g.Open,
		&g.RoleQueue,
		&g.GroupSettings,
		&g.Players,
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
