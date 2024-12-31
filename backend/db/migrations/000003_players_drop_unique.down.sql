ALTER TABLE players
ADD CONSTRAINT players_name_region_key UNIQUE (name, region);