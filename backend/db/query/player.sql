-- name: UpsertPlayer :one
WITH id_check AS (
    SELECT id FROM Players WHERE id = @id
)
INSERT INTO Players (
    name,
    display_name,
    region,
    platform,
    gamemode,
    roles,
    rank,
    characters,
    voice_chat,
    mic,
    vanguards,
    duelists,
    strategists,
    platforms,
    g_voice_chat,
    g_mic
) VALUES (
    @name,
    @display_name,
    @region,
    @platform,
    @gamemode,
    @roles,
    @rank,
    @characters,
    @voice_chat,
    @mic,
    @vanguards,
    @duelists,
    @strategists,
    @platforms,
    @g_voice_chat,
    @g_mic
)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    region = EXCLUDED.region,
    platform = EXCLUDED.platform,
    gamemode = EXCLUDED.gamemode,
    roles = EXCLUDED.roles,
    rank = EXCLUDED.rank,
    characters = EXCLUDED.characters,
    voice_chat = EXCLUDED.voice_chat,
    mic = EXCLUDED.mic,
    vanguards = EXCLUDED.vanguards,
    duelists = EXCLUDED.duelists,
    strategists = EXCLUDED.strategists,
    platforms = EXCLUDED.platforms,
    g_voice_chat = EXCLUDED.g_voice_chat,
    g_mic = EXCLUDED.g_mic,
    updated_at = NOW()
WHERE 
    (SELECT 1 FROM id_check) IS NULL OR -- no specific id provided
    Players.id = @id -- match provided id
RETURNING id;