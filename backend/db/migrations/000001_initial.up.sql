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
    owner VARCHAR(14), -- this is the player's name
    region CHAR(2) NOT NULL,
    gamemode TEXT NOT NULL,
    
    open BOOLEAN NOT NULL,
    passcode VARCHAR(4) NOT NULL DEFAULT generate_passcode(),

    -- role_queue
    vanguards INTEGER NOT NULL DEFAULT 0,
    duelists INTEGER NOT NULL DEFAULT 0,
    strategists INTEGER NOT NULL DEFAULT 0,

    -- group_settings
    platforms CHAR(2)[],
    voice_chat BOOLEAN,
    mic BOOLEAN,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE Ranks (
    id VARCHAR(4) PRIMARY KEY,  
    name VARCHAR(20) NOT NULL, 
    value INTEGER NOT NULL UNIQUE,
    CONSTRAINT valid_rank_id CHECK (id ~ '^(b[1-3]|s[1-3]|g[1-3]|p[1-3]|d[1-3]|gm[1-3]|e|oa)$')
);

INSERT INTO Ranks (id, name, value) VALUES
    ('b3', 'Bronze III', 0),
    ('b2', 'Bronze II', 1),
    ('b1', 'Bronze I', 2),
    ('s3', 'Silver III', 10),
    ('s2', 'Silver II', 11),
    ('s1', 'Silver I', 12),
    ('g3', 'Gold III', 20),
    ('g2', 'Gold II', 21),
    ('g1', 'Gold I', 22),
    ('p3', 'Platinum III', 30),
    ('p2', 'Platinum II', 31),
    ('p1', 'Platinum I', 32),
    ('d3', 'Diamond III', 40),
    ('d2', 'Diamond II', 41),
    ('d1', 'Diamond I', 42),
    ('gm3', 'Grandmaster III', 50),
    ('gm2', 'Grandmaster II', 51),
    ('gm1', 'Grandmaster I', 52),
    ('e', 'Eternity', 60),
    ('oa', 'One Above All', 70);

CREATE OR REPLACE FUNCTION rank_id_to_value(rank_id VARCHAR) 
RETURNS INTEGER AS $$
    SELECT value FROM Ranks WHERE id = $1;
$$ LANGUAGE SQL IMMUTABLE;

CREATE OR REPLACE FUNCTION rank_value_to_id(rank_value INTEGER) 
RETURNS VARCHAR AS $$
    SELECT id FROM Ranks WHERE value = $1;
$$ LANGUAGE SQL IMMUTABLE;

CREATE TABLE Players (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(14) NOT NULL,
    platform CHAR(2) NOT NULL,
    roles TEXT[] NOT NULL,
    rank INTEGER NOT NULL,
    characters TEXT[] NOT NULL,
    voice_chat BOOLEAN NOT NULL,
    mic BOOLEAN NOT NULL,
    
    -- role_queue
    vanguards INTEGER NOT NULL DEFAULT 0,
    duelists INTEGER NOT NULL DEFAULT 0,
    strategists INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE GroupMembers (
    group_id CHAR(4) NOT NULL REFERENCES Groups(id),
    player_id INTEGER NOT NULL REFERENCES Players(id),
    leader BOOLEAN NOT NULL,
    PRIMARY KEY (group_id, player_id)
);