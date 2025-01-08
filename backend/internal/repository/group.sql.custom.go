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
),
requirements_check AS (
    SELECT g.id AS group_id
    FROM Groups g
    LEFT JOIN group_details gd ON g.id = gd.group_id
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
        -- Player requirements when provided
        AND CASE 
            -- Only apply player requirements if player_id is provided
            WHEN $8::INTEGER IS NULL THEN TRUE
            ELSE (
                -- Platform check
                (
                    ARRAY_LENGTH(g.platforms, 1) IS NULL 
                    OR ARRAY_LENGTH(g.platforms, 1) = 0 
                    OR $9::TEXT IS NULL
                    OR $9::TEXT = ANY(g.platforms)
                )
                -- Role queue check
                AND (
                    (g.vanguards + g.duelists + g.strategists = 0)
                    OR
                    ($10::TEXT IS NULL)
                    OR
                    (
                        ($10::TEXT = 'vanguard' AND (
                            SELECT COUNT(*)
                            FROM group_members_base gmb
                            WHERE gmb.group_id = g.id AND gmb.role = 'vanguard'
                        ) < g.vanguards)
                        OR ($10::TEXT = 'duelist' AND (
                            SELECT COUNT(*)
                            FROM group_members_base gmb
                            WHERE gmb.group_id = g.id AND gmb.role = 'duelist'
                        ) < g.duelists)
                        OR ($10::TEXT = 'strategist' AND (
                            SELECT COUNT(*)
                            FROM group_members_base gmb
                            WHERE gmb.group_id = g.id AND gmb.role = 'strategist'
                        ) < g.strategists)
                    )
                )
                -- Rank check (ensure current group members are within range)
                AND ($11::INTEGER IS NULL OR EXISTS (
                    SELECT 1
                    FROM group_members_base gmb
                    WHERE gmb.group_id = g.id
                    AND ABS(
                        (SELECT rank FROM Players WHERE id = gmb.player_id) - $11::INTEGER
                    ) <= 10
                ))
                -- Voice chat and mic requirements
                AND (NOT g.voice_chat OR $12::BOOLEAN IS NULL OR $12::BOOLEAN)
                AND (NOT g.mic OR $13::BOOLEAN IS NULL OR $13::BOOLEAN)
            )
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
        'platforms', g.platforms,
        'voiceChat', g.voice_chat,
        'mic', g.mic
    ) AS group_settings,
    COALESCE(gd.players, '[]'::jsonb) as players,
    COALESCE(gd.member_count, 0) as size,
    CASE WHEN $7 = true THEN (
        SELECT COUNT(*)
        FROM Groups g2
        WHERE g2.id IN (SELECT group_id FROM requirements_check)
    ) ELSE 0 END as total_count,
    g.last_active_at
FROM Groups g
LEFT JOIN group_details gd ON g.id = gd.group_id
WHERE g.id IN (SELECT group_id FROM requirements_check)
ORDER BY 
    CASE WHEN $4 = 'asc' THEN gd.member_count END ASC,
    CASE WHEN $4 = 'desc' THEN gd.member_count END DESC
LIMIT $5 OFFSET $6;
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
	PlayerID  *int32  `json:"playerId"`
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
		arg.SizeSort,
		arg.Limit,
		arg.Offset,
		arg.Count,
		arg.PlayerID,
		arg.Platform,
		arg.Role,
		arg.RankVal,
		arg.VoiceChat,
		arg.Mic,
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
