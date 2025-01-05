import json
import random
import string

def generate_id(existing_ids):
    while True:
        id = ''.join(random.choices(string.ascii_uppercase, k=4))
        if id not in existing_ids:
            return id

# Constants
PLAYERS = [f"Player {i}" for i in range(1, 500)]
REGIONS = ["na", "eu", "ap", "sa"]
PLATFORMS = ["pc", "ps", "xb"]
ROLES = ["vanguard", "duelist", "strategist"]
CHARACTERS = {
    "vanguard": ["Doctor Strange", "Captain America", "Groot", "Hulk", "Magneto", "Peni Parker", "Thor"],
    "duelist": ["Winter Soldier", "Black Panther", "Black Widow"],
    "strategist": ["Rocket Raccoon", "Mantis", "Luna Snow"]
}
RANKS = {
    "b3": 0, "b2": 1, "b1": 2,
    "s3": 10, "s2": 11, "s1": 12,
    "g3": 20, "g2": 21, "g1": 22,
    "p3": 30, "p2": 31, "p1": 32,
    "d3": 40, "d2": 41, "d1": 4
}

def generate_player(player_id, name):
    user_roles = random.sample(ROLES, k=random.randint(1, 3))
    available_characters = []
    for role in user_roles:
        available_characters.extend(CHARACTERS[role])
    characters = random.sample(available_characters, k=random.randint(1, len(available_characters)))
    
    total = 6
    vanguards = random.randint(0, total)
    remaining = total - vanguards
    duelists = random.randint(0, remaining)
    strategists = remaining - duelists

    return (
        f"    ({player_id}, '{name}', '{random.choice(PLATFORMS)}', "
        f"ARRAY{user_roles}, {random.choice(list(RANKS.values()))}, "
        f"ARRAY{characters}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"{vanguards}, {duelists}, {strategists})"
    )

def generate_group(group_id, owner_name):
    platforms = random.choices(PLATFORMS, k=random.randint(1, 3))
    gamemode = random.choice(["competitive", "quickplay"])
    vanguards = random.randint(1, 3)
    duelists = random.randint(1, 3)
    strategists = random.randint(1, 3)
    
    return (
        f"    ('{group_id}', '{owner_name}', '{random.choice(REGIONS)}', "
        f"'{gamemode}', {str(random.choice([True, False])).lower()}, "
        f"DEFAULT, {vanguards}, {duelists}, {strategists}, "
        f"ARRAY{platforms}, {str(random.choice([True, False])).lower()}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"DEFAULT, DEFAULT, DEFAULT)"
    )

def generate_group_member(group_id, player_id, is_leader):
    return f"    ('{group_id}', {player_id}, {str(is_leader).lower()})"

def generate_sql(num_groups=20):
    players_sql = []
    groups_sql = []
    group_members_sql = []
    existing_group_ids = set()
    used_players = set()
    player_id_map = {} 

    for i, name in enumerate(PLAYERS, 1):
        player_id_map[name] = i
        players_sql.append(generate_player(i, name))

    for _ in range(num_groups):
        group_id = generate_id(existing_group_ids)
        existing_group_ids.add(group_id)
        
        owner = random.choice([p for p in PLAYERS if p not in used_players])
        used_players.add(owner)
        owner_id = player_id_map[owner]
        
        groups_sql.append(generate_group(group_id, owner))
        
        group_members_sql.append(generate_group_member(group_id, owner_id, True))
        
        num_additional_members = random.randint(0, 5)
        available_players = [p for p in PLAYERS if p not in used_players]
        if num_additional_members > 0 and available_players:
            additional_members = random.sample(
                available_players,
                k=min(num_additional_members, len(available_players))
            )
            for member in additional_members:
                used_players.add(member)
                member_id = player_id_map[member]
                group_members_sql.append(
                    generate_group_member(group_id, member_id, False)
                )
        
        if len(used_players) >= len(PLAYERS) * 0.8:
            used_players.clear()

    final_sql = [
        "-- Players",
        "INSERT INTO Players (id, name, platform, roles, rank, characters, voice_chat, mic, vanguards, duelists, strategists) VALUES",
        ",\n".join(players_sql) + ";",
        "\n-- Groups",
        "INSERT INTO Groups (id, owner, region, gamemode, open, passcode, vanguards, duelists, strategists, platforms, voice_chat, mic, created_at, updated_at, last_active_at) VALUES",
        ",\n".join(groups_sql) + ";",
        "\n-- GroupMembers",
        "INSERT INTO GroupMembers (group_id, player_id, leader) VALUES",
        ",\n".join(group_members_sql) + ";"
    ]
    
    return "\n".join(final_sql)

if __name__ == "__main__":
    print("Generating mock SQL data...")
    out = generate_sql(250)
    
    with open("mock_data.sql", "w") as f:
        f.write(out)
        print("Saved output to mock_data.sql")