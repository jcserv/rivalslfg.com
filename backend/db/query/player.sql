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
        COUNT(CASE WHEN p.role = 'vanguard' THEN 1 END) as curr_vanguards,
        COUNT(CASE WHEN p.role = 'duelist' THEN 1 END) as curr_duelists,
        COUNT(CASE WHEN p.role = 'strategist' THEN 1 END) as curr_strategists
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
    -- Platform check
    AND g.platform = @platform
    -- Role queue check (only if enabled)
    AND (
        (g.vanguards + g.duelists + g.strategists = 0)
        OR
        (
            -- Can fill at least one role
            (@role = 'vanguard' AND rc.curr_vanguards < g.vanguards)
            OR (@role = 'duelist' AND rc.curr_duelists < g.duelists)
            OR (@role = 'strategist' AND rc.curr_strategists < g.strategists)
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
        role,
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
        @role,
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
    END as status,
    COALESCE(
        (SELECT player_id FROM group_member_creation),
        0
    )::integer as player_id;

-- name: RemovePlayer :one
WITH group_check AS (
    -- Check if group exists and player is in it
    SELECT *
    FROM GroupMembers gm
    WHERE gm.group_id = @group_id
    AND gm.player_id = @player_id
    LIMIT 1
),
group_size AS (
    -- Get total members in group
    SELECT COUNT(*) as member_count
    FROM GroupMembers
    WHERE group_id = @group_id
),
is_last_member AS (
    -- Check if this is the last member
    SELECT (SELECT member_count FROM group_size) = 1 as is_last
),
next_leader AS (
    -- Find next leader if current leader is leaving and not last member
    SELECT gm.player_id
    FROM GroupMembers gm
    JOIN group_check gc ON gm.group_id = gc.group_id
    WHERE gm.player_id != @player_id
    AND NOT EXISTS (SELECT 1 FROM is_last_member WHERE is_last)
    ORDER BY gm.leader DESC, RANDOM()
    LIMIT 1
),
promote_leader AS (
    -- Promote next leader if current leader is leaving and not last member
    UPDATE Groups
    SET owner = (
        SELECT name
        FROM Players
        WHERE id = (SELECT player_id FROM next_leader)
    )
    WHERE id = @group_id
    AND EXISTS (
        SELECT 1 FROM group_check
        WHERE leader = true
    )
    AND EXISTS (SELECT 1 FROM next_leader)
    RETURNING id
),
promote_member AS (
    -- Update group membership for new leader if not last member
    UPDATE GroupMembers
    SET leader = true
    WHERE group_id = @group_id
    AND player_id = (SELECT player_id FROM next_leader)
    AND EXISTS (
        SELECT 1 FROM group_check
        WHERE leader = true
    )
    RETURNING player_id
),
remove_member AS (
    -- Remove the player from the group
    DELETE FROM GroupMembers
    WHERE group_id = @group_id
    AND player_id = @player_id
    AND EXISTS (SELECT 1 FROM group_check)
    RETURNING group_id
),
delete_empty_group AS (
    -- Delete group if this was the last member
    DELETE FROM Groups g
    WHERE g.id = @group_id
    AND EXISTS (SELECT 1 FROM is_last_member WHERE is_last)
    RETURNING id
)
SELECT 
    CASE
        WHEN NOT EXISTS (SELECT 1 FROM group_check) THEN
            '404'::TEXT  -- Group not found or player not in group
        WHEN EXISTS (SELECT 1 FROM is_last_member WHERE is_last) THEN
            '204'::TEXT  -- Last member left, group will be deleted
        ELSE
            '200'::TEXT  -- Successfully removed player
    END as status,
    COALESCE(
        (SELECT player_id FROM next_leader)::INTEGER,
        0
    ) as new_leader_id;