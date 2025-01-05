-- name: CreateGroup :one
-- The result row will contain group_id text and player_id integer
WITH 
-- First check if this combination already exists
existing_membership AS (
    SELECT 
        group_id::text,
        player_id::integer
    FROM GroupMembers
    WHERE 
        CASE 
            WHEN @group_id != '' AND @player_id != 0 THEN 
                group_id = @group_id AND player_id = @player_id
            WHEN @group_id != '' THEN 
                group_id = @group_id
            WHEN @player_id != 0 THEN 
                player_id = @player_id
            ELSE FALSE
        END
    LIMIT 1
),
-- If no membership exists, we might need to create a new player
new_player AS (
    INSERT INTO Players (
        name,
        platform,
        roles,
        rank,
        characters,
        voice_chat,
        mic
    )
    SELECT 
        @owner,
        @platform,
        @roles,
        @rank_val,
        @characters,
        @voice_chat,
        @mic
    WHERE 
        NOT EXISTS (SELECT 1 FROM existing_membership) AND
        (@player_id = 0 OR NOT EXISTS (SELECT 1 FROM Players WHERE id = @player_id))
    RETURNING id
),
-- Get final player_id (either existing, provided, or new)
final_player AS (
    SELECT 
        CASE
            WHEN EXISTS (SELECT 1 FROM existing_membership) THEN 
                (SELECT player_id FROM existing_membership)
            WHEN @player_id != 0 AND EXISTS (SELECT 1 FROM Players WHERE id = @player_id) THEN 
                @player_id
            ELSE 
                (SELECT id FROM new_player)
        END as player_id
),
-- If no membership exists, we might need to create a new group
new_group AS (
    INSERT INTO Groups (
        owner,
        region,
        gamemode,
        open,
        vanguards,
        duelists,
        strategists,
        platforms,
        voice_chat,
        mic
    )
    SELECT
        @owner,
        @region,
        @gamemode,
        @open,
        @vanguards,
        @duelists,
        @strategists,
        @platforms,
        @group_voice_chat,
        @group_mic
    WHERE 
        NOT EXISTS (SELECT 1 FROM existing_membership) AND
        (@group_id = '' OR NOT EXISTS (SELECT 1 FROM Groups WHERE id = @group_id))
    RETURNING id::text
),
-- Get final group_id (either existing, provided, or new)
final_group AS (
    SELECT 
        CASE
            WHEN EXISTS (SELECT 1 FROM existing_membership) THEN 
                (SELECT group_id FROM existing_membership)
            WHEN @group_id != '' AND EXISTS (SELECT 1 FROM Groups WHERE id = @group_id) THEN 
                @group_id
            ELSE 
                (SELECT id FROM new_group)
        END as group_id
),
-- Create the membership if it doesn't exist
new_membership AS (
    INSERT INTO GroupMembers (
        group_id,
        player_id,
        leader
    )
    SELECT 
        fg.group_id,
        fp.player_id,
        true
    FROM final_group fg, final_player fp
    WHERE NOT EXISTS (SELECT 1 FROM existing_membership)
    RETURNING group_id::text, player_id::integer
)
-- Return either the existing or new membership
SELECT group_id::text, player_id::integer
FROM (
    SELECT * FROM existing_membership
    UNION ALL
    SELECT * FROM new_membership
) results
LIMIT 1;