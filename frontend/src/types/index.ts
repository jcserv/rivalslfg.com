import { toTitleCase } from "@/lib/utils";

export enum Region {
  NA = "North America",
  EU = "Europe",
  ME = "Middle East",
  AP = "Asia Pacific",
  SA = "South America",
}

const Regions = {
  na: Region.NA,
  eu: Region.EU,
  me: Region.ME,
  ap: Region.AP,
  sa: Region.SA,
} as const;

export function getRegion(region: string): Region {
  return (
    Object.entries(Regions).find((entry) => entry[0] === region)?.[1] ??
    Region.NA
  );
}

export enum Gamemode {
  Competitive = "competitive",
  Quickplay = "quickplay",
}

export enum Platform {
  PC = "PC",
  PS = "PlayStation",
  XB = "Xbox",
}

const Platforms = {
  pc: {
    emoji: "🖥️",
    label: Platform.PC,
  },
  ps: {
    emoji: "🎮",
    label: Platform.PS,
  },
  xb: {
    emoji: "❎",
    label: Platform.XB,
  },
};

export function getPlatform(platform: string): string {
  const platformObj = Object.entries(Platforms).find(
    (entry) => entry[0] === platform
  )?.[1];
  return `${platformObj?.emoji} ${platformObj?.label}`;
}

export type Role = "vanguard" | "duelist" | "strategist";

export enum Rank {
  b3 = "Bronze III",
  b2 = "Bronze II",
  b1 = "Bronze I",
  s3 = "Silver III",
  s2 = "Silver II",
  s1 = "Silver I",
  g3 = "Gold III",
  g2 = "Gold II",
  g1 = "Gold I",
  p3 = "Platinum III",
  p2 = "Platinum II",
  p1 = "Platinum I",
  d3 = "Diamond III",
  d2 = "Diamond II",
  d1 = "Diamond I",
  gm3 = "Grandmaster III",
  gm2 = "Grandmaster II",
  gm1 = "Grandmaster I",
  e = "Eternity",
  oa = "One Above All",
}

const Ranks = {
  b3: Rank.b3,
  b2: Rank.b2,
  b1: Rank.b1,
  s3: Rank.s3,
  s2: Rank.s2,
  s1: Rank.s1,
  g3: Rank.g3,
  g2: Rank.g2,
  g1: Rank.g1,
  p3: Rank.p3,
  p2: Rank.p2,
  p1: Rank.p1,
  d3: Rank.d3,
  d2: Rank.d2,
  d1: Rank.d1,
  gm3: Rank.gm3,
  gm2: Rank.gm2,
  gm1: Rank.gm1,
  e: Rank.e,
  oa: Rank.oa,
} as const;

export function getRank(rank: string): Rank {
  return (
    Object.entries(Ranks).find((entry) => entry[0] === rank)?.[1] ?? Rank.b3
  );
}

export type Profile = {
  name: string;
  region: Region;
  platform: Platform;
  gamemode: Gamemode;
  roles: Role[];
  rank: Rank;
  characters: string[];
  roleQueue?: RoleQueue;
};

export type RoleQueue = {
  vanguards: number;
  duelists: number;
  strategists: number;
};

export type Group = {
  region: string;
  gamemode: string;
  roleQueue?: RoleQueue;
  players: Player[];
};

export type Player = {
  name: string;
  leader?: boolean;
  rank: string;
  roles: string[];
  characters: string[];
  platform: string;
};

export type TeamUp = {
  name: string;
  requirements: {
    allOf: string[];
    oneOf?: string[];
  };
  description: string;
  seasonBonus: {
    stat: string;
    modifier: string;
    target: string;
    value: number;
    unit: string;
  };
};

function formatStat(stat: TeamUp["seasonBonus"]["stat"]) {
  switch (stat) {
    case "max-health":
      return "Max Health";
    case "healing":
      return "Healing";
    case "damage":
      return "Damage";
    default:
      return stat;
  }
}

export function formatSeasonBonus(seasonBonus: TeamUp["seasonBonus"]) {
  return `${seasonBonus.target}: ${seasonBonus.unit === "+" ? "+" : ""}${seasonBonus.value}${
    seasonBonus.unit === "%" ? "%" : ""
  } ${formatStat(seasonBonus.stat)} ${toTitleCase(seasonBonus.modifier)}`;
}
