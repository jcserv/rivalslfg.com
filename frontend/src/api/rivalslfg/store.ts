import { Store } from "@tanstack/store";

import { Group } from "@/types";
import { queryClient } from "@/routes/__root";

interface RivalsLFGStore {
  groups: Group[];
}

const initialState: RivalsLFGStore = {
  groups: [],
};

export const rivalslfgStore = new Store<RivalsLFGStore>(initialState);

export const rivalsStoreKeys = {
  all: ["groups"] as const,
  group: (id: string) => ["group", id] as const,
};

export const rivalsStoreActions = {
  setGroups(groups: Group[]) {
    rivalslfgStore.setState((prev) => ({ ...prev, groups }));
  },
  setGroup(group: Group) {
    rivalslfgStore.setState((prev) => ({
      ...prev,
      groups: [...prev.groups, group],
    }));
  },
  clearStore() {
    rivalslfgStore.setState(() => initialState);
    queryClient.removeQueries({ queryKey: rivalsStoreKeys.all });
  },
};
