-- name: CreateGroupWithOwner :one
WITH player_info AS (
    SELECT Players.id AS player_id, Players.display_name
    FROM Players
    WHERE Players.id = $1 -- player_id
),
new_group AS (
    INSERT INTO Groups (
        community_id,
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
        $2,                  -- community_id
        player_info.display_name,  -- owner
        $3,                -- region
        $4,                -- gamemode
        $5,                -- open
        $6,                -- vanguards
        $7,                -- duelists
        $8,                -- strategists
        $9,                -- platforms
        $10,                -- voice_chat
        $11                -- mic
    FROM player_info
    RETURNING id
)
INSERT INTO GroupMembers (group_id, player_id, leader)
SELECT 
    new_group.id,
    player_info.player_id,
    TRUE
FROM new_group, player_info
RETURNING group_id;