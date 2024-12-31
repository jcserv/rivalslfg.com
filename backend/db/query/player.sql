-- name: CreatePlayer :one
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
    $1, -- name
    $2, -- display_name
    $3, -- region
    $4, -- platform
    $5, -- gamemode
    $6, -- roles
    $7, -- rank
    $8, -- characters
    $9, -- voice_chat
    $10, -- mic

    $11, -- vanguards
    $12, -- duelists
    $13, -- strategists

    $14, -- platforms
    $15, -- g_voice_chat
    $16  -- g_mic
) RETURNING id;