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

-- name: JoinGroup :one
WITH group_status AS (
    SELECT 
        CASE
            WHEN NOT EXISTS (SELECT 1 FROM Groups WHERE g.id = $1) THEN 404
            WHEN NOT g.open AND g.passcode != $2 THEN 403
            ELSE 200
        END as status,
        players
    FROM Groups g
    WHERE id = $1
)
UPDATE Groups g
SET 
    players = CASE 
        WHEN (SELECT status FROM group_status) = 200 
        THEN jsonb_insert(COALESCE(players, '[]'::jsonb), '{-1}', sqlc.arg(player)::jsonb)
        ELSE players
    END,
    last_active_at = NOW(),
    updated_at = NOW()
WHERE g.id = $1
RETURNING (SELECT status FROM group_status) as result;