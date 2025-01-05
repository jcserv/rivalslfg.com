import random
import json
import string
from typing import Optional, Union, List, Dict, Any

class MockLobbyGenerator:
    def __init__(self):
        self.regions = ['na', 'eu', 'me', 'ap', 'sa']
        self.gamemodes = ['competitive', 'quickplay']
        self.ranks = [
            'b3', 'b2', 'b1',  # Bronze
            's3', 's2', 's1',  # Silver
            'g3', 'g2', 'g1',  # Gold
            'p3', 'p2', 'p1',  # Platinum
            'd3', 'd2', 'd1',  # Diamond
            'gm3', 'gm2', 'gm1'  # Grandmaster
            'e', 'oa'  # Eternity, One Above All
        ]
        self.roles = ['vanguard', 'duelist', 'strategist']
        self.characters = [
            'Doctor Strange', 'Luna Snow', 'Mantis', 'Rocket Raccoon',
            'Winter Soldier', 'Spider-Man', 'Black Panther', 'Storm',
            'Magneto', 'Star-Lord', 'Iron Man', 'Captain America'
        ]
        self.platforms = ['pc', 'ps', 'xb']
        self.usernames = [
            'Skelzore', 'imphungky', 'xZestence', 'Scynthesia',
            'StarGuardian', 'NightHawk', 'CosmicRider', 'ThunderBolt',
            'FrostByte', 'PhoenixRise', 'ShadowStrike', 'LightBearer',
            'DarkMatter', 'VoidWalker', 'StormChaser', 'MindBender'
        ]
        self.used_ids = set()

    def generate_unique_id(self) -> str:
        """Generate a unique 4-letter ID."""
        while True:
            new_id = ''.join(random.choices(string.ascii_uppercase, k=4))
            if new_id not in self.used_ids:
                self.used_ids.add(new_id)
                return new_id
            
    def generate_role_queue(self) -> Optional[Dict[str, int]]:
        """Generate role queue settings with a probability of 80%."""
        if random.random() < 0.8:  # 80% chance to have role queue
            # Generate random distribution that sums to 6
            while True:
                vanguards = random.randint(1, 4)
                duelists = random.randint(1, 4)
                strategists = 6 - vanguards - duelists
                if strategists >= 1:  # Ensure at least 1 of each role
                    return {
                        "vanguards": vanguards,
                        "duelists": duelists,
                        "strategists": strategists
                    }
        return None
    
    def generate_group_settings(self) -> Dict[str, Any]:
        """Generate group settings."""
        # Random number of platforms (can be empty)
        num_platforms = random.randint(0, len(self.platforms))
        selected_platforms = random.sample(self.platforms, num_platforms)
        
        return {
            "platforms": selected_platforms,
            "voiceChat": random.choice([True, False]),
            "mic": random.choice([True, False])
        }

    def generate_player(self, is_leader: bool = False) -> Dict[str, Any]:
        """
        Generate a single player's data.
        
        Args:
            is_leader (bool): Whether this player should be the leader
        """
        role = random.choice(self.roles)
        return {
            "name": random.choice(self.usernames),
            "leader": is_leader,
            "rank": random.choice(self.ranks),
            "roles": [role],
            "characters": random.sample(self.characters, random.randint(1, 3)),
            "platform": random.choice(self.platforms)
        }

    def generate_lobby(self) -> Dict[str, Any]:
        """Generate a single lobby's data."""
        player_count = random.randint(2, 6)
        
        # Generate leader first
        leader = self.generate_player(is_leader=True)
        
        # Generate other players (non-leaders)
        other_players = [self.generate_player(is_leader=False) for _ in range(player_count - 1)]
        
        # Combine all players with leader first
        all_players = [leader] + other_players
        
        lobby = {
            "id": self.generate_unique_id(),
            "name": f"{leader['name']}'s Group",
            "region": random.choice(self.regions),
            "gamemode": random.choice(self.gamemodes),
            "open": random.choice([True, False]),
            "groupSettings": self.generate_group_settings(),
            "players": all_players
        }

        role_queue = self.generate_role_queue()
        if role_queue:
            lobby["roleQueue"] = role_queue
        
        return lobby

    def generate(self, count: int = 1) -> Union[Dict[str, Any], List[Dict[str, Any]]]:
        """
        Generate mock lobby data.
        
        Args:
            count (int): Number of lobbies to generate. Defaults to 1.
            
        Returns:
            If count == 1, returns a single lobby dictionary.
            If count > 1, returns a list of lobby dictionaries.
        """
        self.used_ids = set()

        if count == 1:
            return self.generate_lobby()
        return [self.generate_lobby() for _ in range(count)]

    def generate_and_save(self, count: int, filename: str):
        """
        Generate mock lobby data and save it to a JSON file.
        
        Args:
            count (int): Number of lobbies to generate
            filename (str): Name of the file to save the data to
        """
        lobbies = self.generate(count)
        
        # Ensure the filename has .json extension
        if not filename.endswith('.json'):
            filename += '.json'
            
        # Save to file with pretty printing (saving just the array)
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(lobbies, f, indent=2)
        
        print(f"Generated {count} lobbies and saved to {filename}")


if __name__ == "__main__":
    # Create generator instance
    generator = MockLobbyGenerator()
    
    # Generate 250 lobbies and save them to a file
    generator.generate_and_save(250, 'mock_lobbies.json')