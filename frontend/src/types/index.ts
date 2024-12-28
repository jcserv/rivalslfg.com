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
  characterPrefs: string[];
};
