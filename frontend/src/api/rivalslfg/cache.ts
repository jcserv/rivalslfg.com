import { rivalslfgStore, rivalsStoreActions } from "@/api/rivalslfg/store";
import { rivalslfgAPIClient } from "@/routes/__root";
import { getGroupFromProfile, Group, Profile } from "@/types";
import { StatusCode, StatusCodes } from "@/types/http";

export const upsertGroup = async (
  profile: Profile,
  id: string,
): Promise<string> => {
  const groupId = await rivalslfgAPIClient.upsertGroup(profile, id);
  const newGroup = getGroupFromProfile(profile, groupId);
  rivalsStoreActions.setAuthedGroup(groupId);

  // TODO: It'd be nice to have an upsertGroup action in the store so that this is just one action
  if (id) {
    rivalsStoreActions.replaceGroup(newGroup);
  } else {
    rivalsStoreActions.setGroup(newGroup);
  }

  return groupId;
};

export const fetchGroups = async (): Promise<Group[]> => {
  const groups = await rivalslfgAPIClient.getGroups();
  rivalsStoreActions.setGroups(groups);
  return groups;
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
  rivalsStoreActions.setGroup(group);
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
