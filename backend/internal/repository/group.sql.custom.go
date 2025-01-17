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
        LOWER(p.role) as role,
        p.rank as rank_val,
        rank_value_to_id(p.rank) as rank_id,
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
                'rank', rank_id,
                'characters', characters,
                'voiceChat', voice_chat,
                'mic', mic
            )
        ) as players,
        MIN(rank_val) as min_rank,
        MAX(rank_val) as max_rank,
        COUNT(CASE WHEN role = 'vanguard' THEN 1 END) as curr_vanguards,
        COUNT(CASE WHEN role = 'duelist' THEN 1 END) as curr_duelists,
        COUNT(CASE WHEN role = 'strategist' THEN 1 END) as curr_strategists
    FROM group_members_base
    GROUP BY group_id
),
requirements_check AS (
    SELECT g.id AS group_id
    FROM Groups g
    JOIN group_details gd ON g.id = gd.group_id
    WHERE
        -- Base requirements
        ($1 = '' OR g.region = $1)
        AND ($2 = '' OR g.gamemode = $2)
        AND ($3 = '' OR
            CASE 
                WHEN LOWER($3) = 'true' THEN g.open = true
                WHEN LOWER($3) = 'false' THEN g.open = false
                ELSE TRUE
            END
        )
        -- Player requirements check
        AND CASE 
            -- If rank value is provided, use it as a trigger for all player requirements
            WHEN $8::INTEGER IS NOT NULL THEN (
                -- Platform check
                g.platform = $4
                -- Role queue check
                AND (
                    g.vanguards + g.duelists + g.strategists = 0
                    OR (
                        CASE $5::TEXT
                            WHEN 'vanguard' THEN gd.curr_vanguards < g.vanguards 
                            WHEN 'duelist' THEN gd.curr_duelists < g.duelists
                            WHEN 'strategist' THEN gd.curr_strategists < g.strategists
                            ELSE FALSE
                        END
                    )
                )
                -- Rank check
                AND (
                    -- Allow Bronze-Gold players to group with each other
                    ($8::INTEGER BETWEEN 0 AND 22 AND gd.min_rank BETWEEN 0 AND 22)
                    OR (
                        ABS(gd.min_rank - $8::INTEGER) <= 10 
                        AND ABS(gd.max_rank - $8::INTEGER) <= 10
                    )
                )
                -- Voice chat and mic
                AND (NOT g.voice_chat OR $6::BOOLEAN)
                AND (NOT g.mic OR $7::BOOLEAN)
            )
            ELSE TRUE
        END
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
        'platform', g.platform,
        'voiceChat', g.voice_chat,
        'mic', g.mic
    ) AS group_settings,
    COALESCE(gd.players, '[]'::jsonb) as players,
    COALESCE(gd.member_count, 0) as size,
    CASE WHEN $9 = true THEN (
        SELECT COUNT(*)
        FROM Groups g2
        WHERE g2.id IN (SELECT group_id FROM requirements_check)
    ) ELSE 0 END as total_count,
    g.last_active_at
FROM Groups g
JOIN group_details gd ON g.id = gd.group_id
WHERE g.id IN (SELECT group_id FROM requirements_check)
ORDER BY 
    CASE WHEN $10 = 'asc' THEN gd.member_count END ASC,
    CASE WHEN $10 = 'desc' THEN gd.member_count END DESC
LIMIT $11 OFFSET $12;
`

type GetGroupsParams struct {
	RegionFilter   string `json:"regionFilter"`
	GamemodeFilter string `json:"gamemodeFilter"`

	// OpenFilter is a string to account for when we don't want to filter
	OpenFilter string `json:"openFilter"`
	SizeSort   string `json:"sizeSort"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Count      bool   `json:"count"`

	// Player requirements (all optional)
	Platform  *string `json:"platform"`
	Role      *string `json:"role"`
	RankVal   *int32  `json:"rankVal"`
	VoiceChat *bool   `json:"voiceChat"`
	Mic       *bool   `json:"mic"`
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
	rows, err := q.db.Query(ctx, getGroups,
		arg.RegionFilter,
		arg.GamemodeFilter,
		arg.OpenFilter,
		arg.Platform,
		arg.Role,
		arg.VoiceChat,
		arg.Mic,
		arg.RankVal,
		arg.Count,
		arg.SizeSort,
		arg.Limit,
		arg.Offset,
	)
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
        'platform', g.platform,
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
