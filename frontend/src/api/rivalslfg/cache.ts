import {
  rivalslfgStore,
  rivalsStoreActions,
  StatusCode,
  StatusCodes,
} from "@/api";
import { rivalslfgAPIClient } from "@/routes/__root";
import {
  CreateGroupResponse,
  getGroupFromProfile,
  Group,
  PaginatedQueryFnResponse,
  Profile,
  QueryParams,
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
  requesterId: number,
): Promise<StatusCode> => {
  const result = await rivalslfgAPIClient.removePlayer(groupId, requesterId);
  // if (result === StatusCodes.OK) {
  //   rivalsStoreActions.removePlayerFromGroup(playerId);
  // }
  return result;
};
