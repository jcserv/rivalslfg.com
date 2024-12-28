enum Platform {
  PC = "PC",
  PlayStation = "PlayStation",
  Xbox = "Xbox",
}

export type Profile = {
  region: string;
  platform: Platform;
  gamemode: string;
  roles: string[];
  rank: string;
  characters: string[];
  roleQueue?: RoleQueue;
};

export type RoleQueue = {
  vanguards: number;
  duelists: number;
  strategists: number;
};
