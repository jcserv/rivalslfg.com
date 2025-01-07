import {
  HTTPError,
  rivalslfgStore,
  rivalsStoreActions,
  StatusCode,
  StatusCodes,
} from "@/api";
import { rivalslfgAPIClient } from "@/routes/__root";
import {
  CreateGroupResponse,
  Gamemode,
  getGroupFromProfile,
  Group,
  PaginatedQueryFnResponse,
  Profile,
  QueryParams,
  Region,
} from "@/types";

export const createGroup = async (
  profile: Profile,
): Promise<CreateGroupResponse> => {
  const { groupId, playerId } = await rivalslfgAPIClient.createGroup(profile);
  const newGroup = getGroupFromProfile(profile, groupId);

  rivalsStoreActions.setAuthedGroup(groupId);
  rivalsStoreActions.upsertGroup(newGroup);

  return { groupId, playerId };
};

export const fetchGroups = async (
  query?: QueryParams,
): PaginatedQueryFnResponse<Group> => {
  const { groups, totalCount } = await rivalslfgAPIClient.getGroups(query);
  rivalsStoreActions.setGroups(groups);
  return { data: groups, totalCount };
};

export const fetchGroup = async (id: string): Promise<Group | undefined> => {
  const cached = rivalslfgStore.state.groups.find((group) => group.id === id);
  if (cached) {
    return cached;
  }

  try {
    const group = await rivalslfgAPIClient.getGroup(id);
    if (!group) {
      return undefined;
    }
    rivalsStoreActions.upsertGroup(group);
    return group;
  } catch (error) {
    if (!(error instanceof HTTPError)) {
      return undefined;
    }
    if (!(error.statusCode === StatusCodes.Forbidden)) {
      throw error;
    }
    // For forbidden groups, return a minimal group object
    return {
      id,
      open: false,
      // Add other required fields with placeholder/default values
      name: "Private Group",
      owner: "",
      ownerId: 0,
      region: "" as Region,
      gamemode: "" as Gamemode,
      players: [],
      groupSettings: {
        platforms: [],
        voiceChat: false,
        mic: false,
      },
    };
  }
};

export const joinGroup = async (
  groupId: string,
  player: Profile,
  passcode: string,
): Promise<StatusCode> => {
  const result = await rivalslfgAPIClient.joinGroup(groupId, player, passcode);
  if (result === StatusCodes.OK) {
    rivalsStoreActions.setAuthedGroup(groupId);
  }
  return result;
};

export const removePlayer = async (
  groupId: string,
  playerId: number,
): Promise<StatusCode> => {
  const result = await rivalslfgAPIClient.removePlayer(groupId, playerId);
  if (result === StatusCodes.NoContent) {
    rivalsStoreActions.removeGroup(groupId);
    return result;
  }
  if (result === StatusCodes.OK) {
    rivalsStoreActions.removePlayerFromGroup(groupId, playerId);
  }
  return result;
};
