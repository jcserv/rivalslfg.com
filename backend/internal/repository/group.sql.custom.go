package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const getGroups = `
WITH group_members_base AS (
    SELECT 
        gm.group_id,
        p.id as player_id,
        p.name,
        gm.leader,
        p.platform,
        p.role,
        rank_value_to_id(p.rank) as rank,
        p.characters,
        p.voice_chat,
        p.mic
    FROM GroupMembers gm
    JOIN Players p ON p.id = gm.player_id
),
group_details AS (
    SELECT 
        group_id,
        COUNT(*) as member_count,
        MAX(CASE WHEN leader = true THEN player_id END) as owner_id,
        jsonb_agg(
            jsonb_build_object(
                'id', player_id,
                'name', name,
                'leader', leader,
                'platform', platform,
                'role', role,
                'rank', rank,
                'characters', characters,
                'voiceChat', voice_chat,
                'mic', mic
            )
        ) as players
    FROM group_members_base
    GROUP BY group_id
)

SELECT 
    g.id,
    community_id,
	COALESCE(gd.owner_id, 0) as owner_id,
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
    COALESCE(gd.players, '[]'::jsonb) as players,
    COALESCE(gd.member_count, 0) as size,
    CASE WHEN $7 = true THEN (
        SELECT COUNT(*)
        FROM Groups g2
        WHERE ($1 = '' OR g2.region = $1)
          AND ($2 = '' OR g2.gamemode = $2)
          AND ($3 = '' OR
            CASE 
                WHEN LOWER($3) = 'true' THEN g2.open = true
                WHEN LOWER($3) = 'false' THEN g2.open = false
                ELSE TRUE
            END
          )
    ) ELSE 0 END as total_count,
    g.last_active_at
FROM Groups g
LEFT JOIN group_details gd ON g.id = gd.group_id
WHERE ($1 = '' OR g.region = $1)
  AND ($2 = '' OR g.gamemode = $2)
  AND ($3 = '' OR
    CASE 
        WHEN LOWER($3) = 'true' THEN g.open = true
        WHEN LOWER($3) = 'false' THEN g.open = false
        ELSE TRUE
    END
  )
ORDER BY 
    CASE WHEN $4 = 'asc' THEN gd.member_count END ASC,
    CASE WHEN $4 = 'desc' THEN gd.member_count END DESC
LIMIT $5 OFFSET $6;
`

type GetGroupsParams struct {
	RegionFilter   string `json:"regionFilter"`
	GamemodeFilter string `json:"gamemodeFilter"`
	OpenFilter     string `json:"openFilter"` // openFilter is a string to account for when we don't want to filter
	SizeSort       string `json:"sizeSort"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Count          bool   `json:"count"`
}

type GetGroupsRow struct {
	GroupWithPlayers
	TotalCount int32 `json:"totalCount"`
}

func (g *GetGroupsRow) ToGroupWithPlayers() GroupWithPlayers {
	return GroupWithPlayers{
		GroupDTO: GroupDTO{
			ID:            g.ID,
			CommunityID:   g.CommunityID,
			OwnerID:       g.OwnerID,
			Owner:         g.Owner,
			Region:        g.Region,
			Gamemode:      g.Gamemode,
			Open:          g.Open,
			Passcode:      g.Passcode,
			RoleQueue:     g.RoleQueue,
			GroupSettings: g.GroupSettings,
			LastActiveAt:  g.LastActiveAt,
		},
		Name:    g.Name,
		Size:    g.Size,
		Players: g.Players,
	}
}

type GetGroupsResult struct {
	Groups     []GroupWithPlayers `json:"groups"`
	TotalCount int32              `json:"totalCount"`
}

func (q *Queries) GetGroups(ctx context.Context, arg GetGroupsParams) (*GetGroupsResult, error) {
	rows, err := q.db.Query(ctx, getGroups, arg.RegionFilter, arg.GamemodeFilter, arg.OpenFilter,
		arg.SizeSort, arg.Limit, arg.Offset, arg.Count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &GetGroupsResult{
		Groups:     make([]GroupWithPlayers, 0),
		TotalCount: 0,
	}

	for rows.Next() {
		var g GetGroupsRow
		if err := rows.Scan(
			&g.ID,
			&g.CommunityID,
			&g.OwnerID,
			&g.Owner,
			&g.Region,
			&g.Gamemode,
			&g.Open,
			&g.RoleQueue,
			&g.GroupSettings,
			&g.Players,
			&g.Size,
			&g.TotalCount,
			&g.LastActiveAt,
		); err != nil {
			return nil, err
		}
		g.Name = fmt.Sprintf("%s's Group", g.Owner)
		result.TotalCount = g.TotalCount
		result.Groups = append(result.Groups, g.ToGroupWithPlayers())
	}
	return result, nil
}

const getGroupByID = `
WITH group_members AS (
    SELECT 
        gm.group_id,
        gm.leader,
        p.id as player_id,
        p.name,
        p.platform,
        p.role,
        rank_value_to_id(p.rank) as rank,
        p.characters,
        p.voice_chat,
        p.mic
    FROM GroupMembers gm
    JOIN Players p ON p.id = gm.player_id
    WHERE gm.group_id = $1
)
SELECT 
    g.id,
    community_id,
	(
        SELECT gm2.player_id 
        FROM group_members gm2 
        WHERE gm2.leader = true 
        LIMIT 1
    ) as owner_id,
    owner,
    g.region,
    g.gamemode,
    open,
    g.passcode,
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
        jsonb_agg(
            jsonb_build_object(
                'id', gm.player_id,
                'name', gm.name,
                'leader', gm.leader,
                'platform', gm.platform,
                'role', gm.role,
                'rank', gm.rank,
                'characters', gm.characters,
                'voiceChat', gm.voice_chat,
                'mic', gm.mic
            )
        ),
        '[]'::jsonb
    ) as players,
    COUNT(gm.player_id) AS size,
    g.last_active_at
FROM Groups g
LEFT JOIN group_members gm ON g.id = gm.group_id
WHERE g.id = $1
GROUP BY g.id`

func (q *Queries) GetGroupByID(ctx context.Context, id string) (*GroupWithPlayers, error) {
	var g GroupWithPlayers
	err := q.db.QueryRow(ctx, getGroupByID, id).Scan(
		&g.ID,
		&g.CommunityID,
		&g.OwnerID,
		&g.Owner,
		&g.Region,
		&g.Gamemode,
		&g.Open,
		&g.Passcode,
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
