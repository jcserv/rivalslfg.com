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
  Group,
  JoinGroupResponse,
  PaginatedQueryFnResponse,
  Platform,
  Profile,
  QueryParams,
  Region,
} from "@/types";

export const createGroup = async (
  profile: Profile,
): Promise<CreateGroupResponse> => {
  const { groupId, playerId } = await rivalslfgAPIClient.createGroup(profile);

  rivalsStoreActions.setAuthedGroup(groupId);

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
      name: "Private Group",
      owner: "",
      ownerId: 0,
      region: "" as Region,
      gamemode: "" as Gamemode,
      players: [],
      groupSettings: {
        platform: Platform.PC,
        voiceChat: false,
        mic: false,
      },
    };
  }
};

export const patchGroup = async (
  id: string,
  open: boolean,
): Promise<StatusCode> => {
  const cached = rivalslfgStore.state.groups.find((group) => group.id === id);
  if (!cached) return StatusCodes.NotFound;

  const response = await rivalslfgAPIClient.patchGroup(id, { open });
  if (response === StatusCodes.NoContent) {
    rivalsStoreActions.upsertGroup({
      ...cached,
      open,
    });
  }
  return response;
};

export const deleteGroup = async (id: string): Promise<StatusCode> => {
  const response = await rivalslfgAPIClient.deleteGroup(id);
  if (response === StatusCodes.NoContent) {
    rivalsStoreActions.removeAuthedGroup(id);
    rivalsStoreActions.removeGroup(id);
  }
  return response;
};

export const joinGroup = async (
  groupId: string,
  player: Profile,
  passcode: string,
): Promise<JoinGroupResponse> => {
  const response = await rivalslfgAPIClient.joinGroup(
    groupId,
    player,
    passcode,
  );
  if (response.status === StatusCodes.OK) {
    rivalsStoreActions.setAuthedGroup(groupId);
  }
  return response;
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
  return result;
};
