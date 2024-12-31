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

export const gamemodeEmojis: Record<Gamemode, string> = {
  [Gamemode.Competitive]: "👑",
  [Gamemode.Quickplay]: "⚡",
};

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
    (entry) => entry[0] === platform,
  )?.[1];
  return `${platformObj?.emoji} ${platformObj?.label}`;
}

export type Role = "vanguard" | "duelist" | "strategist";
export const Roles = ["vanguard", "duelist", "strategist"] as const;

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

function rankToKey(rank: Rank): RankKey {
  return rank.toString() as RankKey;
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

type RankKey = keyof typeof Ranks;

export function getRank(rank: string): Rank {
  return (
    Object.entries(Ranks).find((entry) => entry[0] === rank)?.[1] ?? Rank.b3
  );
}

const RankVals: Record<RankKey, number> = {
  b3: 0,
  b2: 1,
  b1: 2,
  s3: 10,
  s2: 11,
  s1: 12,
  g3: 20,
  g2: 21,
  g1: 22,
  p3: 30,
  p2: 31,
  p1: 32,
  d3: 40,
  d2: 41,
  d1: 42,
  gm3: 50,
  gm2: 51,
  gm1: 52,
  e: 60,
  oa: 70,
} as const;

export const getRankFromRankVal = (rankVal: number): Rank => {
  const rankKey =
    Object.entries(RankVals).find((entry) => entry[1] === rankVal)?.[0] ?? "b3";
  return getRank(rankKey);
};

export function isAdjacentRank(
  userRank: RankKey,
  comparisonRank: number,
): boolean {
  return Math.abs(RankVals[userRank] - comparisonRank) <= 10;
}

export function canJoinGroup(userRank: RankKey, groupRanks: RankKey[]) {
  return (
    groupRanks.reduce((acc, rank) => {
      return acc + (RankVals[rank] - RankVals[userRank]);
    }, 0) <= 10
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
  voiceChat: boolean;
  mic: boolean;
  roleQueue?: RoleQueue;
  groupSettings?: GroupSettings;
};

export type RoleQueue = {
  vanguards: number;
  duelists: number;
  strategists: number;
};

export type GroupSettings = {
  platforms: Platform[];
  voiceChat: boolean;
  mic: boolean;
};

export type Group = {
  id: string;
  name: string;
  open: boolean;
  region: Region;
  gamemode: Gamemode;
  players: Player[];
  groupSettings: GroupSettings;
  roleQueue?: RoleQueue;
};

type GroupInfo = {
  minRank: number;
  maxRank: number;
  currVanguards: number;
  currDuelists: number;
  currStrategists: number;
  currCharacters: Set<string>;
};

export function getGroupInfo(group: Group): GroupInfo {
  return group.players.reduce(
    (acc, player) => {
      const rankKey = player.rank as RankKey;
      acc.minRank = Math.min(acc.minRank, RankVals[rankKey]);
      acc.maxRank = Math.max(acc.maxRank, RankVals[rankKey]);
      acc.currVanguards += player.roles.includes("vanguard") ? 1 : 0;
      acc.currDuelists += player.roles.includes("duelist") ? 1 : 0;
      acc.currStrategists += player.roles.includes("strategist") ? 1 : 0;
      acc.currCharacters = acc.currCharacters.union(new Set(player.characters));
      return acc;
    },
    {
      minRank: RankVals["oa"],
      maxRank: RankVals["b1"],
      currVanguards: 0,
      currDuelists: 0,
      currStrategists: 0,
      currCharacters: new Set<string>(),
    },
  );
}

export type GroupRequirements = {
  minRank: number;
  maxRank: number;
  requestedRoles: {
    vanguards: {
      curr: number;
      max: number;
    };
    duelists: {
      curr: number;
      max: number;
    };
    strategists: {
      curr: number;
      max: number;
    };
  };
  voiceChat: boolean;
  mic: boolean;
  platforms: Platform[];
};

export function getRequirements(group: Group): GroupRequirements {
  const { currVanguards, currDuelists, currStrategists, minRank, maxRank } =
    getGroupInfo(group);
  return {
    minRank,
    maxRank,
    requestedRoles: {
      vanguards: {
        curr: currVanguards,
        max: group.roleQueue?.vanguards ?? 0,
      },
      duelists: {
        curr: currDuelists,
        max: group.roleQueue?.duelists ?? 0,
      },
      strategists: {
        curr: currStrategists,
        max: group.roleQueue?.strategists ?? 0,
      },
    },
    platforms: group.groupSettings.platforms,
    mic: group.groupSettings.mic,
    voiceChat: group.groupSettings.voiceChat,
  };
}

export function areRequirementsMet(
  group: Group,
  requirements: GroupRequirements,
  profile: Profile | null,
): boolean {
  const { minRank, maxRank, mic, voiceChat, platforms } = requirements;
  if (!profile) return false;
  if (!Object.keys(profile).length) return false;
  const { rank, gamemode, region, platform } = profile;
  const rankKey = rankToKey(rank);

  const basicRequirementsMet =
    gamemode === group.gamemode &&
    region === group.region &&
    (platforms.length > 0 ? platforms.includes(platform) : true) &&
    (mic ? profile.mic : true) &&
    (voiceChat ? profile.voiceChat : true);

  if (!basicRequirementsMet) return false;
  if (!group.roleQueue) return true;

  const canFill = profile.roles.some((role) => {
    switch (role.toLowerCase()) {
      case Roles[0]: {
        const vanguardSpots = requirements.requestedRoles.vanguards;
        return vanguardSpots.curr < vanguardSpots.max;
      }
      case Roles[1]: {
        const duelistSpots = requirements.requestedRoles.duelists;
        return duelistSpots.curr < duelistSpots.max;
      }
      case Roles[2]: {
        const strategistSpots = requirements.requestedRoles.strategists;
        return strategistSpots.curr < strategistSpots.max;
      }
      default:
        return false;
    }
  });

  return (
    isAdjacentRank(rankKey, minRank) &&
    isAdjacentRank(rankKey, maxRank) &&
    canFill
  );
}

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
