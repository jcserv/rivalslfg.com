import { Store } from "@tanstack/react-store";

import { queryClient } from "@/routes/__root";
import { Group } from "@/types";

interface RivalsLFGStore {
  authedGroups: Map<string, boolean>;
  groups: Group[];
}

const initialState: RivalsLFGStore = {
  authedGroups: new Map<string, boolean>(),
  groups: [],
};

export const rivalslfgStore = new Store<RivalsLFGStore>(initialState);

export const rivalsStoreKeys = {
  all: ["authedGroups", "groups"] as const,
  groups: ["groups"] as const,
  group: (id: string) => ["group", id] as const,
};

export const rivalsStoreActions = {
  setAuthedGroup(groupId: string) {
    rivalslfgStore.setState((prev) => ({
      ...prev,
      authedGroups: new Map([...prev.authedGroups, [groupId, true]]),
    }));
  },
  setGroups(groups: Group[]) {
    rivalslfgStore.setState((prev) => ({ ...prev, groups }));
  },
  upsertGroup(group: Group) {
    rivalslfgStore.setState((prev) => ({
      ...prev,
      groups: prev.groups.some((g) => g.id === group.id)
        ? prev.groups.map((g) => (g.id === group.id ? group : g))
        : [group, ...prev.groups],
    }));
  },
  clearStore() {
    rivalslfgStore.setState(() => initialState);
    queryClient.removeQueries({ queryKey: rivalsStoreKeys.all });
  },
};
