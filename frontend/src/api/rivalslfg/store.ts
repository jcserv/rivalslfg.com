import { Store } from "@tanstack/react-store";

import { Group } from "@/types";
import { queryClient } from "@/routes/__root";

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
  setGroup(group: Group) {
    rivalslfgStore.setState((prev) => ({
      ...prev,
      groups: [group, ...prev.groups],
    }));
  },
  replaceGroup(group: Group) {
    rivalslfgStore.setState((prev) => ({
      ...prev,
      groups: prev.groups.map((g) => (g.id === group.id ? group : g)),
    }));
  },
  clearStore() {
    rivalslfgStore.setState(() => initialState);
    queryClient.removeQueries({ queryKey: rivalsStoreKeys.all });
  },
};
