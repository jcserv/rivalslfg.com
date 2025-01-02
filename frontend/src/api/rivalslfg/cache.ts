import {
  PaginationParams,
  rivalslfgStore,
  rivalsStoreActions,
  StatusCode,
  StatusCodes,
} from "@/api";
import { rivalslfgAPIClient } from "@/routes/__root";
import {
  getGroupFromProfile,
  Group,
  PaginatedQueryFnResponse,
  Profile,
} from "@/types";

export const upsertGroup = async (
  profile: Profile,
  id: string,
): Promise<string> => {
  const groupId = await rivalslfgAPIClient.upsertGroup(profile, id);
  const newGroup = getGroupFromProfile(profile, groupId);

  rivalsStoreActions.setAuthedGroup(groupId);
  rivalsStoreActions.upsertGroup(newGroup);

  return groupId;
};

export const fetchGroups = async (
  pagination?: PaginationParams,
): PaginatedQueryFnResponse<Group> => {
  const { groups, totalCount } = await rivalslfgAPIClient.getGroups(pagination);
  rivalsStoreActions.setGroups(groups);
  return { data: groups, totalCount };
};

export const fetchGroup = async (id: string): Promise<Group | undefined> => {
  const cached = rivalslfgStore.state.groups.find((group) => group.id === id);
  if (cached) {
    return cached;
  }

  const group = await rivalslfgAPIClient.getGroup(id);
  if (!group) {
    return undefined;
  }
  rivalsStoreActions.upsertGroup(group);
  return group;
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
  requesterName: string,
  playerName: string,
): Promise<StatusCode> => {
  const result = await rivalslfgAPIClient.removePlayer(
    groupId,
    playerId,
    requesterName,
    playerName,
  );
  // if (result === StatusCodes.OK) {
  //   rivalsStoreActions.removePlayerFromGroup(playerId);
  // }
  return result;
};
