-- name: UpsertGroup :one
WITH id_check AS (
    SELECT id FROM Groups WHERE id = @id
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
        WHEN @id IS NULL OR @id = '' THEN generate_group_id()
        ELSE @id
    END,
    @owner,
    @region,
    @gamemode,
    COALESCE(@players, '[]'::jsonb),
    @open,
    @vanguards,
    @duelists,
    @strategists,
    @platforms,
    @voice_chat,
    @mic,
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
    Groups.id = @id -- match provided id
RETURNING id;

-- name: CheckCanJoinGroup :one
SELECT 
    CASE
        WHEN NOT EXISTS (SELECT 1 FROM Groups WHERE g.id = @id) THEN 404
        WHEN NOT g.open AND g.passcode != @passcode THEN 403
        WHEN EXISTS (
            SELECT 1 FROM jsonb_array_elements(g.players) AS p
            WHERE p->>'name' = @player_name::text
        ) THEN 200
        ELSE 202
    END as status
FROM Groups g
WHERE g.id = @id;

-- name: JoinGroup :exec
UPDATE Groups g
SET 
    players = jsonb_insert(COALESCE(players, '[]'::jsonb), '{-1}', @player::jsonb),
    last_active_at = NOW(),
    updated_at = NOW()
WHERE g.id = @id;

-- name: RemovePlayerFromGroup :one
WITH group_check AS (
    SELECT 
        CASE
            WHEN NOT EXISTS (SELECT 1 FROM Groups WHERE g.id = @id) THEN 404
            WHEN @requester_name != g.owner 
                AND @requester_name != @player_name THEN 403
            WHEN NOT EXISTS (
                SELECT 1 FROM jsonb_array_elements(players) AS p
                WHERE p->>'name' = @player_name::text
            ) THEN 404
            ELSE 200
        END as status,
        players
    FROM Groups g
    WHERE id = @id
),
player_update as (
    UPDATE Groups g
    SET 
        players = COALESCE(
            (
                SELECT jsonb_agg(value)
                FROM jsonb_array_elements(g.players) AS p
                WHERE p->>'name' != @player_name::text
            ),
            '[]'::jsonb
        ),
        last_active_at = NOW(),
        updated_at = NOW()
    WHERE g.id = @id
        AND (SELECT status FROM group_check) = 200
)
SELECT status FROM group_check;