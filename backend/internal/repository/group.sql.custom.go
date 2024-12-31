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
    COALESCE(
        json_agg(
            jsonb_build_object(
                'id', p.id,
                'name', p.name,
                'displayName', p.display_name,
                'region', p.region,
                'platform', p.platform,
                'gamemode', p.gamemode,
                'roles', p.roles,
                'rank', p.rank,
                'characters', p.characters,
                'voiceChat', p.voice_chat,
                'mic', p.mic,
				'roleQueue', jsonb_build_object(
					'vanguards', p.vanguards,
					'duelists', p.duelists,
					'strategists', p.strategists
				),
				'groupSettings', jsonb_build_object(
					'platforms', p.platforms,
					'voiceChat', p.g_voice_chat,
					'mic', p.g_mic
				)
            )
        ), '[]'
    ) as players
FROM Groups g
LEFT JOIN GroupMembers gm ON g.id = gm.group_id
LEFT JOIN Players p ON p.id = gm.player_id
GROUP BY g.id`

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
	COALESCE(
		json_agg(
			jsonb_build_object(
				'id', p.id,
                'name', p.name,
                'displayName', p.display_name,
                'region', p.region,
                'platform', p.platform,
                'gamemode', p.gamemode,
                'roles', p.roles,
                'rank', p.rank,
                'characters', p.characters,
                'voiceChat', p.voice_chat,
                'mic', p.mic,
				'roleQueue', jsonb_build_object(
					'vanguards', p.vanguards,
					'duelists', p.duelists,
					'strategists', p.strategists
				),
				'groupSettings', jsonb_build_object(
					'platforms', p.platforms,
					'voiceChat', p.g_voice_chat,
					'mic', p.g_mic
				)
			)
		) FILTER (WHERE p.id IS NOT NULL), 
		'[]'
	) as players
FROM Groups g
LEFT JOIN GroupMembers gm ON g.id = gm.group_id
LEFT JOIN Players p ON p.id = gm.player_id
WHERE g.id = $1
GROUP BY g.id`

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
