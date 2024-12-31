package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const findAllGroups = `
SELECT 
    g.*,
    COALESCE(
        json_agg(
            jsonb_build_object(
                'id', p.id,
                'name', p.name,
                'display_name', p.display_name,
                'region', p.region,
                'platform', p.platform,
                'gamemode', p.gamemode,
                'roles', p.roles,
                'rank', p.rank,
                'characters', p.characters,
                'p_voice_chat', p.p_voice_chat,
                'p_mic', p.p_mic,
                'vanguards', p.vanguards,
                'duelists', p.duelists,
                'strategists', p.strategists,
                'platforms', p.platforms,
                'g_voice_chat', p.g_voice_chat,
                'g_mic', p.g_mic,
                'created_at', p.created_at,
                'updated_at', p.updated_at
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
			&g.Passcode,
			&g.Vanguards,
			&g.Duelists,
			&g.Strategists,
			&g.Platforms,
			&g.VoiceChat,
			&g.Mic,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.Players,
		); err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, nil
}

const findGroupByID = `
SELECT 
	g.*,
	COALESCE(
		json_agg(
			jsonb_build_object(
				'id', p.id,
				'name', p.name,
				'display_name', p.display_name,
				'region', p.region,
				'platform', p.platform,
				'gamemode', p.gamemode,
				'roles', p.roles,
				'rank', p.rank,
				'characters', p.characters,
				'p_voice_chat', p.p_voice_chat,
				'p_mic', p.p_mic,
				'vanguards', p.vanguards,
				'duelists', p.duelists,
				'strategists', p.strategists,
				'platforms', p.platforms,
				'g_voice_chat', p.g_voice_chat,
				'g_mic', p.g_mic,
				'created_at', p.created_at,
				'updated_at', p.updated_at
			)
		) FILTER (WHERE p.id IS NOT NULL), 
		'[]'
	) as players
FROM Groups g
LEFT JOIN GroupMembers gm ON g.id = gm.group_id
LEFT JOIN Players p ON p.id = gm.player_id
WHERE g.id = $1
GROUP BY g.id, g.community_id, g.owner, g.region, g.gamemode, g.open, 
			g.passcode, g.vanguards, g.duelists, g.strategists, g.platforms, 
			g.voice_chat, g.mic, g.created_at, g.updated_at`

func (q *Queries) FindGroupByID(ctx context.Context, id string) (*GroupWithPlayers, error) {
	var g GroupWithPlayers
	err := q.db.QueryRow(ctx, findGroupByID, id).Scan(
		&g.ID,
		&g.CommunityID,
		&g.Owner,
		&g.Region,
		&g.Gamemode,
		&g.Open,
		&g.Passcode,
		&g.Vanguards,
		&g.Duelists,
		&g.Strategists,
		&g.Platforms,
		&g.VoiceChat,
		&g.Mic,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.Players,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}
