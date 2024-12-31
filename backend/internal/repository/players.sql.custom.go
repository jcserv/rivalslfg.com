package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const findPlayer = `
SELECT 
    id,
    name,
    display_name,
    region,
    platform,
    gamemode,
    roles,
    rank,
    characters,
    voice_chat,
    mic,
    jsonb_build_object(
        'vanguards', vanguards,
        'duelists', duelists,
        'strategists', strategists
    ) AS role_queue,
    jsonb_build_object(
        'platforms', platforms,
        'g_voice_chat', voice_chat,
        'g_mic', mic
    ) AS group_settings
FROM Players
WHERE (id = $1 OR name = $2)
`

func (q *Queries) FindPlayer(ctx context.Context, id int32, name string) (*Player, error) {
	var p Player
	err := q.db.QueryRow(ctx, findPlayer, id, name).Scan(
		&p.ID,
		&p.Name,
		&p.DisplayName,
		&p.Region,
		&p.Platform,
		&p.Gamemode,
		&p.Roles,
		&p.Rank,
		&p.Characters,
		&p.VoiceChat,
		&p.Mic,
		&p.RoleQueue,
		&p.GroupSettings,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
