import { Group, Profile } from "./types";

export interface PaginatedGroupsResponse {
  groups: Group[];
  pageCount: number;
  totalCount: number;
}

export type CreateGroup = {
  owner: string;
  region: string;
  gamemode: string;
  open: boolean;

  platform: string;
  role: string;
  rankId: string;
  characters: string[];
  voiceChat: boolean;
  mic: boolean;

  vanguards: number;
  duelists: number;
  strategists: number;

  platforms: string[];
  groupVoiceChat: boolean;
  groupMic: boolean;
};

export function getCreateGroupFromProfile(profile: Profile): CreateGroup {
  return {
    owner: profile.name,
    region: profile.region,
    gamemode: profile.gamemode,
    open: true,

    platform: profile.platform,
    role: profile.role,
    rankId: profile.rank,
    characters: profile.characters,
    voiceChat: profile.voiceChat,
    mic: profile.mic,

    vanguards: profile.roleQueue?.vanguards ?? 0,
    duelists: profile.roleQueue?.duelists ?? 0,
    strategists: profile.roleQueue?.strategists ?? 0,

    platforms: [],
    groupVoiceChat: false,
    groupMic: false,
  };
}

export type CreateGroupResponse = {
  groupId: string;
  playerId: number;
};
