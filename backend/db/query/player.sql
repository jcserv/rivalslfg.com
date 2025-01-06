-- name: JoinGroup :one
WITH 
-- First check if player is already in a group
player_check AS (
    SELECT 1 FROM GroupMembers gm
    WHERE gm.player_id = @player_id
    LIMIT 1
),

-- Get current role counts for the group
role_counts AS (
    SELECT 
        COUNT(CASE WHEN 'vanguard' = ANY(p.roles) THEN 1 END) as curr_vanguards,
        COUNT(CASE WHEN 'duelist' = ANY(p.roles) THEN 1 END) as curr_duelists,
        COUNT(CASE WHEN 'strategist' = ANY(p.roles) THEN 1 END) as curr_strategists
    FROM GroupMembers gm
    JOIN Players p ON p.id = gm.player_id
    WHERE gm.group_id = @group_id
),

-- Check all requirements in a single query
valid_group AS (
    SELECT g.id
    FROM Groups g, role_counts rc
    WHERE g.id = @group_id
    AND g.gamemode = @gamemode
    AND g.region = @region
    -- Platform check (only if platforms specified)
    AND (
        ARRAY_LENGTH(g.platforms, 1) IS NULL 
        OR ARRAY_LENGTH(g.platforms, 1) = 0 
        OR @platform::TEXT = ANY(g.platforms)
    )
    -- Role queue check (only if enabled)
    AND (
        (g.vanguards + g.duelists + g.strategists = 0)
        OR
        (
            -- Can fill at least one role
            ('vanguard' = ANY(@roles) AND rc.curr_vanguards < g.vanguards)
            OR ('duelist' = ANY(@roles) AND rc.curr_duelists < g.duelists)
            OR ('strategist' = ANY(@roles) AND rc.curr_strategists < g.strategists)
        )
    )
    -- Rank check
    AND EXISTS (
        SELECT 1
        FROM Players p2
        WHERE p2.id IN (SELECT player_id FROM GroupMembers WHERE group_id = g.id)
        AND ABS(p2.rank - @rank_val) <= 10
    )
    -- If group is not open, check if passcode is correct
    AND (
        g.open 
        OR (NOT g.open AND g.passcode = @passcode)
    )
    LIMIT 1
),

-- Insert player if they don't exist and group is valid
player_creation AS (
    INSERT INTO Players (
        name,
        platform,
        roles,
        rank,
        characters,
        voice_chat,
        mic,
        vanguards,
        duelists,
        strategists
    )
    SELECT 
        @name,
        @platform::TEXT,
        @roles,
        @rank_val,
        @characters,
        @voice_chat,
        @mic,
        @vanguards,
        @duelists,
        @strategists
    WHERE 
        NOT EXISTS (SELECT 1 FROM player_check)
        AND EXISTS (SELECT 1 FROM valid_group)
        AND NOT EXISTS (SELECT 1 FROM Players WHERE id = @player_id)
    RETURNING id
),

-- Create group membership if everything valid
group_member_creation AS (
    INSERT INTO GroupMembers (
        group_id,
        player_id,
        leader
    )
    SELECT 
        @group_id,
        COALESCE(
            (SELECT id FROM player_creation),
            @player_id
        ),
        false
    WHERE EXISTS (SELECT 1 FROM valid_group)
    AND NOT EXISTS (SELECT 1 FROM player_check)
    RETURNING player_id
)

-- Return status code
SELECT 
    CASE
        WHEN EXISTS (SELECT 1 FROM player_check) THEN '400a'
        WHEN NOT EXISTS (SELECT 1 FROM Groups g WHERE g.id = @group_id) THEN '404'
        WHEN EXISTS (SELECT 1 FROM group_member_creation) THEN '200'
        WHEN EXISTS (
            SELECT 1 FROM Groups g 
            WHERE g.id = @id 
            AND NOT g.open 
            AND g.passcode != @passcode
        ) THEN '403'
        WHEN NOT EXISTS (SELECT 1 FROM valid_group) THEN '400e'
        ELSE '500'
    END as status;