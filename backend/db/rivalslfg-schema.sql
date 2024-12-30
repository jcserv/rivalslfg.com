CREATE TYPE Region AS ENUM (
    'na',
    'eu',
    'me',
    'ap',
    'sa'
);

CREATE TYPE Gamemode AS ENUM (
    'competitive',
    'quickplay'
);

CREATE TYPE Platform AS ENUM (
    'pc',
    'ps',
    'xb'
);

CREATE TYPE Role AS ENUM (
    'vanguard',
    'duelist',
    'strategist'
);

CREATE TYPE RankName AS ENUM (
    "Bronze III",
    "Bronze II",
    "Bronze I",
    "Silver III",
    "Silver II",
    "Silver I",
    "Gold III",
    "Gold II",
    "Gold I",
    "Platinum III",
    "Platinum II",
    "Platinum I",
    "Diamond III",
    "Diamond II",
    "Diamond I",
    "Grandmaster III",
    "Grandmaster II",
    "Grandmaster I",
    "Eternity",
    "One Above All"
);

CREATE TYPE RankID AS ENUM (
    "b3", -- 0
    "b2", -- 1
    "b1", -- 2
    "s3", -- 10
    "s2", -- 11
    "s1", -- 12
    "g3", -- 20
    "g2", -- 21
    "g1", -- 22
    "p3", -- 30
    "p2", -- 31
    "p1", -- 32
    "d3", -- 40
    "d2", -- 41
    "d1", -- 42
    "gm3", -- 50
    "gm2", -- 51
    "gm1", -- 52
    "e", -- 60
    "oa" -- 70
);

CREATE TYPE Rank AS (
    id RankID,
    name RankName,
    val integer
);