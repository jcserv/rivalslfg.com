import random
import string

def generate_id(existing_ids):
    while True:
        id = ''.join(random.choices(string.ascii_uppercase, k=4))
        if id not in existing_ids:
            return id

PLAYERS = [f"Player {i}" for i in range(1, 500)]
REGIONS = ["na", "eu", "ap", "sa"]
PLATFORMS = ["pc", "co"]
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
    "d3": 40, "d2": 41, "d1": 42,
    "gm3": 50, "gm2": 51, "gm1": 52,
    "c3": 60, "c2": 61, "c1": 62,
    "e": 70, "oa": 80
}

RANK_VALS_TO_ID = {v: k for k, v in RANKS.items()}

def is_adjacent_rank(rank_val1, rank_val2):
    """Check if two rank values are adjacent (within 10 units of each other)."""
    return abs(rank_val1 - rank_val2) <= 10

def get_adjacent_rank(target_rank_val):
    """Get a random rank value that is adjacent to the target rank value."""
    possible_ranks = [
        rank_val for rank_val in RANKS.values()
        if is_adjacent_rank(rank_val, target_rank_val)
    ]
    return random.choice(possible_ranks)

def generate_player(name, rank_val=None):
    role = random.choice(ROLES)
    available_characters = CHARACTERS[role]
    characters = random.sample(available_characters, k=random.randint(1, len(available_characters)))
    
    # If no rank_val provided, generate a random one
    if rank_val is None:
        rank_val = random.choice(list(RANKS.values()))
    
    total = 6
    vanguards = random.randint(0, total)
    remaining = total - vanguards
    duelists = random.randint(0, remaining)
    strategists = remaining - duelists

    return (
        f"    ('{name}', '{random.choice(PLATFORMS)}', "
        f"'{role}', {rank_val}, "
        f"ARRAY{characters}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"{vanguards}, {duelists}, {strategists})"
    )

def generate_group(group_id, owner_name):
    platform = random.choice(PLATFORMS)
    gamemode = random.choice(["competitive", "quickplay"])
    vanguards = random.randint(1, 3)
    duelists = random.randint(1, 3)
    strategists = random.randint(1, 3)
    
    return (
        f"    ('{group_id}', '{owner_name}', '{random.choice(REGIONS)}', "
        f"'{gamemode}', {str(random.choice([True, False])).lower()}, "
        f"DEFAULT, {vanguards}, {duelists}, {strategists}, "
        f"'{platform}', {str(random.choice([True, False])).lower()}, "
        f"{str(random.choice([True, False])).lower()}, "
        f"DEFAULT, DEFAULT, DEFAULT)"
    )

def generate_group_member(group_id, player_name, is_leader):
    return f"    ('{group_id}', (SELECT id FROM Players WHERE name = '{player_name}'), {str(is_leader).lower()})"

def generate_sql(num_groups=20):
    players_sql = []
    groups_sql = []
    group_members_sql = []
    existing_group_ids = set()
    used_players = set()
    player_ranks = {}  # Track player ranks for group formation

    # First generate all players with random ranks
    for name in PLAYERS:
        rank_val = random.choice(list(RANKS.values()))
        players_sql.append(generate_player(name, rank_val))
        player_ranks[name] = rank_val

    for _ in range(num_groups):
        group_id = generate_id(existing_group_ids)
        existing_group_ids.add(group_id)
        
        # Select owner from available players
        available_owners = [p for p in PLAYERS if p not in used_players]
        if not available_owners:
            used_players.clear()
            available_owners = PLAYERS
            
        owner = random.choice(available_owners)
        owner_rank = player_ranks[owner]
        used_players.add(owner)
        
        groups_sql.append(generate_group(group_id, owner))
        group_members_sql.append(generate_group_member(group_id, owner, True))
        
        # Add additional members with adjacent ranks
        num_additional_members = random.randint(0, 5)
        if num_additional_members > 0:
            # Filter available players by rank adjacency
            available_players = [
                p for p in PLAYERS 
                if p not in used_players and is_adjacent_rank(player_ranks[p], owner_rank)
            ]
            
            if available_players:
                additional_members = random.sample(
                    available_players,
                    k=min(num_additional_members, len(available_players))
                )
                for member in additional_members:
                    used_players.add(member)
                    group_members_sql.append(
                        generate_group_member(group_id, member, False)
                    )

        if len(used_players) >= len(PLAYERS) * 0.8:
            used_players.clear()

    final_sql = [
        "-- Players",
        "INSERT INTO Players (name, platform, role, rank, characters, voice_chat, mic, vanguards, duelists, strategists) VALUES",
        ",\n".join(players_sql) + ";",
        "\n-- Groups",
        "INSERT INTO Groups (id, owner, region, gamemode, open, passcode, vanguards, duelists, strategists, platform, voice_chat, mic, created_at, updated_at, last_active_at) VALUES",
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