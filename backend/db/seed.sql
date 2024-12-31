INSERT INTO Community (name, description, link) VALUES
    ('Rivals LFG', 'A site that helps Marvel Rivals players find groups to play with', 'https://rivalslfg.com');

INSERT INTO Groups (id, owner, region, gamemode, open, vanguards, duelists, strategists, platforms, voice_chat, mic) VALUES
    ('AAAA', 'Skelzore', 'na', 'competitive', true, 2, 2, 2, ARRAY['pc'], false, false),
    ('AAAB', 'imphungky', 'na', 'competitive', false, 2, 2, 2, ARRAY['pc'], false, false);

INSERT INTO Players (name, display_name, region, platform, gamemode, roles, rank, characters, p_voice_chat, p_mic) VALUES
    ('skelzore', 'Skelzore', 'na', 'pc', 'competitive', ARRAY['strategist'], 'p3', ARRAY['Mantis'], false, false),
    ('imphungky', 'imphungky', 'na', 'xb', 'competitive', ARRAY['vanguard'], 'd3', ARRAY['Doctor Strange'], false, false),
    ('xzestence', 'xZestence', 'na', 'ps', 'quickplay', ARRAY['strategist'], 'p1', ARRAY['Rocket Raccoon', 'Luna Snow'], true, false),
    ('scynthesia', 'Scynthesia', 'na', 'pc', 'quickplay', ARRAY['duelist'], 'd3', ARRAY['Winter Solider'], false, false);

INSERT INTO GroupMembers (group_id, player_id, leader) VALUES
    ('AAAA', 1, true),
    ('AAAA', 2, false),
    ('AAAA', 3, false),
    ('AAAA', 4, false);