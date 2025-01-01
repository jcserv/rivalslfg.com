// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: group.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkCanJoinGroup = `-- name: CheckCanJoinGroup :one
SELECT 
    CASE
        WHEN NOT EXISTS (SELECT 1 FROM Groups WHERE g.id = $1) THEN 404
        WHEN NOT g.open AND g.passcode != $2 THEN 403
        WHEN EXISTS (
            SELECT 1 FROM jsonb_array_elements(g.players) AS p
            WHERE p->>'name' = $3::text
        ) THEN 200
        ELSE 202
    END as status
FROM Groups g
WHERE g.id = $1
`

type CheckCanJoinGroupParams struct {
	ID         string `json:"id"`
	Passcode   string `json:"passcode"`
	PlayerName string `json:"player_name"`
}

func (q *Queries) CheckCanJoinGroup(ctx context.Context, arg CheckCanJoinGroupParams) (int32, error) {
	row := q.db.QueryRow(ctx, checkCanJoinGroup, arg.ID, arg.Passcode, arg.PlayerName)
	var status int32
	err := row.Scan(&status)
	return status, err
}

const joinGroup = `-- name: JoinGroup :exec
UPDATE Groups g
SET 
    players = jsonb_insert(COALESCE(players, '[]'::jsonb), '{-1}', $1::jsonb),
    last_active_at = NOW(),
    updated_at = NOW()
WHERE g.id = $2
`

type JoinGroupParams struct {
	Player []byte `json:"player"`
	ID     string `json:"id"`
}

func (q *Queries) JoinGroup(ctx context.Context, arg JoinGroupParams) error {
	_, err := q.db.Exec(ctx, joinGroup, arg.Player, arg.ID)
	return err
}

const removePlayerFromGroup = `-- name: RemovePlayerFromGroup :one
WITH group_check AS (
    SELECT 
        CASE
            WHEN NOT EXISTS (SELECT 1 FROM Groups WHERE g.id = $1) THEN 404
            WHEN $2 != g.owner 
                AND $2 != $3 THEN 403
            WHEN NOT EXISTS (
                SELECT 1 FROM jsonb_array_elements(players) AS p
                WHERE p->>'name' = $3::text
            ) THEN 404
            ELSE 200
        END as status,
        players
    FROM Groups g
    WHERE id = $1
),
player_update as (
    UPDATE Groups g
    SET 
        players = COALESCE(
            (
                SELECT jsonb_agg(value)
                FROM jsonb_array_elements(g.players) AS p
                WHERE p->>'name' != $3::text
            ),
            '[]'::jsonb
        ),
        last_active_at = NOW(),
        updated_at = NOW()
    WHERE g.id = $1
        AND (SELECT status FROM group_check) = 200
)
SELECT status FROM group_check
`

type RemovePlayerFromGroupParams struct {
	ID            string      `json:"id"`
	RequesterName string      `json:"requester_name"`
	PlayerName    interface{} `json:"player_name"`
}

func (q *Queries) RemovePlayerFromGroup(ctx context.Context, arg RemovePlayerFromGroupParams) (int32, error) {
	row := q.db.QueryRow(ctx, removePlayerFromGroup, arg.ID, arg.RequesterName, arg.PlayerName)
	var status int32
	err := row.Scan(&status)
	return status, err
}

const upsertGroup = `-- name: UpsertGroup :one
WITH id_check AS (
    SELECT id FROM Groups WHERE id = $1
)
INSERT INTO Groups (
    id,
    owner,
    region,
    gamemode,
    players,
    open,
    vanguards,
    duelists,
    strategists,
    platforms,
    voice_chat,
    mic,
    last_active_at
) VALUES (
    CASE 
        WHEN $1 IS NULL OR $1 = '' THEN generate_group_id()
        ELSE $1
    END,
    $2,
    $3,
    $4,
    COALESCE($5, '[]'::jsonb),
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    owner = EXCLUDED.owner,
    region = EXCLUDED.region,
    gamemode = EXCLUDED.gamemode,
    players = EXCLUDED.players,
    open = EXCLUDED.open,
    vanguards = EXCLUDED.vanguards,
    duelists = EXCLUDED.duelists,
    strategists = EXCLUDED.strategists,
    platforms = EXCLUDED.platforms,
    voice_chat = EXCLUDED.voice_chat,
    mic = EXCLUDED.mic,
    last_active_at = NOW(),
    updated_at = NOW()
WHERE 
    (SELECT 1 FROM id_check) IS NULL OR -- no specific id provided
    Groups.id = $1 -- match provided id
RETURNING id
`

type UpsertGroupParams struct {
	ID          interface{} `json:"id"`
	Owner       string      `json:"owner"`
	Region      string      `json:"region"`
	Gamemode    string      `json:"gamemode"`
	Players     interface{} `json:"players"`
	Open        bool        `json:"open"`
	Vanguards   pgtype.Int4 `json:"vanguards"`
	Duelists    pgtype.Int4 `json:"duelists"`
	Strategists pgtype.Int4 `json:"strategists"`
	Platforms   []string    `json:"platforms"`
	VoiceChat   pgtype.Bool `json:"voice_chat"`
	Mic         pgtype.Bool `json:"mic"`
}

func (q *Queries) UpsertGroup(ctx context.Context, arg UpsertGroupParams) (string, error) {
	row := q.db.QueryRow(ctx, upsertGroup,
		arg.ID,
		arg.Owner,
		arg.Region,
		arg.Gamemode,
		arg.Players,
		arg.Open,
		arg.Vanguards,
		arg.Duelists,
		arg.Strategists,
		arg.Platforms,
		arg.VoiceChat,
		arg.Mic,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}
