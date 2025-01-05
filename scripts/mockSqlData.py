import json
import random
import string

def generate_id(existing_ids):
    while True:
        id = ''.join(random.choices(string.ascii_uppercase, k=4))
        if id not in existing_ids:
            return id

OWNERS = [f"Player {i}" for i in range(1, 500)]
REGIONS = ["na", "eu", "ap", "sa"]
PLATFORMS = ["pc", "ps", "xb"]
ROLES = ["vanguard", "duelist", "strategist"]
CHARACTERS = {
    "vanguard": ["Doctor Strange", "Captain America", "Groot", "Hulk", "Magneto", "Peni Parker", "Thor"],
    "duelist": ["Winter Soldier", "Black Panther", "Black Widow"],
    "strategist": ["Rocket Raccoon", "Mantis", "Luna Snow"]
}
RANKS = ["b3", "b2", "b1", "s3", "s2", "s1", "g3", "g2", "g1", "p3", "p2", "p1", "d3", "d2", "d1"]

def generate_player(name, is_leader=False):
    user_roles = random.choices(ROLES, k=random.randint(1, 3))
    available_characters = []
    for role in user_roles:
        available_characters.extend(CHARACTERS[role])
    characters = random.sample(available_characters, k=random.randint(0, len(available_characters)))

    return {
        "id": random.randint(1, 1000),
        "name": name,
        "leader": is_leader,
        "platform": random.choice(PLATFORMS),
        "roles": user_roles,
        "rank": random.choice(RANKS),
        "characters": characters,
        "voiceChat": random.choice([True, False]),
        "mic": random.choice([True, False])
    }

def generate_group(group_id, used_players):
    owner = random.choice([p for p in OWNERS if p not in used_players])
    used_players.add(owner)
    
    num_players = random.randint(1, 6)
    players = [generate_player(owner, True)]
    
    available_players = [p for p in OWNERS if p not in used_players]
    if num_players > 1 and available_players:
        additional_players = random.sample(available_players, k=min(num_players-1, len(available_players)))
        for player in additional_players:
            used_players.add(player)
            players.append(generate_player(player))

    platforms = random.choices(PLATFORMS, k=random.randint(1, 3))
    gamemode = random.choice(["competitive", "quickplay"])
    vanguards = random.randint(1, 3)
    duelists = random.randint(1, 3)
    strategists = random.randint(1, 3)
    
    return (
        f"    ('{group_id}', '{owner}', '{random.choice(REGIONS)}', '{gamemode}', "
        f"{str(random.choice([True, False])).lower()}, {vanguards}, {duelists}, {strategists}, "
        f"ARRAY{platforms}, {str(random.choice([True, False])).lower()}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"'{json.dumps(players)}'::json)"
    )

def generate_sql(num_groups=20):
    sql = "INSERT INTO Groups (id, owner, region, gamemode, open, vanguards, duelists, strategists, platforms, voice_chat, mic, players) VALUES\n"
    values = []
    existing_ids = set()
    used_players = set()
    
    for _ in range(num_groups):
        group_id = generate_id(existing_ids)
        existing_ids.add(group_id)
        values.append(generate_group(group_id, used_players))
        if len(used_players) >= len(OWNERS):
            used_players.clear()
    
    sql += ",\n".join(values) + ";"
    return sql

if __name__ == "__main__":
    print("Generating mock SQL data...")
    out = generate_sql(250)

    with open("mock_lobbies.sql", "w") as f:
        f.write(out)
        print("Saved output to mock_lobbies.sql")