import { rivalslfgStore, rivalsStoreActions } from "@/api/rivalslfg/store";
import { rivalslfgAPIClient } from "@/routes/__root";
import { Group } from "@/types";

export const fetchGroups = async (): Promise<Group[]> => {
  const cached = await rivalslfgStore.state.groups;
  if (cached.length > 0) {
    return cached;
  }

  const groups = await rivalslfgAPIClient.getGroups();
  rivalsStoreActions.setGroups(groups);
  return groups;
};

export const fetchGroup = async (id: string): Promise<Group> => {
  const cached = (await rivalslfgStore.state.groups).find(
    (group) => group.id === id,
  );
  if (cached) {
    return cached;
  }

  const group = await rivalslfgAPIClient.getGroup(id);
  rivalsStoreActions.setGroup(group);
  return group;
};
