-- Types

CREATE TYPE RankName AS ENUM (
    'Bronze III',
    'Bronze II',
    'Bronze I',
    'Silver III',
    'Silver II',
    'Silver I',
    'Gold III',
    'Gold II',
    'Gold I',
    'Platinum III',
    'Platinum II',
    'Platinum I',
    'Diamond III',
    'Diamond II',
    'Diamond I',
    'Grandmaster III',
    'Grandmaster II',
    'Grandmaster I',
    'Eternity',
    'One Above All'
);

CREATE TYPE RankID AS ENUM (
    'b3', -- 0
    'b2', -- 1
    'b1', -- 2
    's3', -- 10
    's2', -- 11
    's1', -- 12
    'g3', -- 20
    'g2', -- 21
    'g1', -- 22
    'p3', -- 30
    'p2', -- 31
    'p1', -- 32
    'd3', -- 40
    'd2', -- 41
    'd1', -- 42
    'gm3', -- 50
    'gm2', -- 51
    'gm1', -- 52
    'e', -- 60
    'oa' -- 70
);

CREATE TYPE Rank AS (
    id RankID,
    name RankName,
    val integer
);

-- Functions

CREATE OR REPLACE FUNCTION generate_group_id() 
RETURNS char(4) AS $$
DECLARE
    chars char[] := ARRAY['A','B','C','D','E','F','G','H','I','J','K','L','M',
                         'N','O','P','Q','R','S','T','U','V','W','X','Y','Z'];
    result char(4) := '';
    i integer := 0;
BEGIN
    -- Generate a random 4-letter string
    WHILE i < 4 LOOP
        result := result || chars[1 + floor(random() * 26)];
        i := i + 1;
    END LOOP;
    
    -- If ID already exists, try again (recursive)
    IF EXISTS (SELECT 1 FROM groups WHERE id = result) THEN
        result := generate_group_id();
    END IF;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION generate_passcode() 
RETURNS char(4) AS $$
DECLARE
    chars char[] := ARRAY['A','B','C','D','E','F','G','H','I','J','K','L','M',
                         'N','O','P','Q','R','S','T','U','V','W','X','Y','Z',
                         '0','1','2','3','4','5','6','7','8','9'];
    result char(4) := '';
    i integer := 0;
BEGIN
    -- Generate a random 4-character string
    WHILE i < 4 LOOP
        result := result || chars[1 + floor(random() * 36)]; -- 26 letters + 10 numbers = 36
        i := i + 1;
    END LOOP;
    
    -- If passcode already exists, try again (recursive)
    IF EXISTS (SELECT 1 FROM groups WHERE passcode = result) THEN
        result := generate_passcode();
    END IF;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Tables

CREATE TABLE Community (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL
);

CREATE TABLE Groups (
    id CHAR(4) PRIMARY KEY DEFAULT generate_group_id(),
    community_id INTEGER NOT NULL REFERENCES Community(id) DEFAULT 1,
    owner VARCHAR(14) NOT NULL, -- this is the player's name
    region CHAR(2) NOT NULL,
    gamemode TEXT NOT NULL,
    players JSONB NOT NULL DEFAULT '[]'::JSONB,
    
    open BOOLEAN NOT NULL,
    passcode VARCHAR(4) NOT NULL DEFAULT generate_passcode(),

    -- role_queue
    vanguards INTEGER,
    duelists INTEGER,
    strategists INTEGER,

    -- group_settings
    platforms CHAR(2)[],
    voice_chat BOOLEAN,
    mic BOOLEAN,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);