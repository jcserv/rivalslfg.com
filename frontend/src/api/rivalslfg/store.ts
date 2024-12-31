import { Store } from "@tanstack/store";

import { Group, Profile } from "@/types";
import { queryClient } from "@/routes/__root";

interface RivalsLFGStore {
  groups: Group[];
  profile: Profile | null;
}

const initialState: RivalsLFGStore = {
  groups: [],
  profile: null,
};

export const rivalslfgStore = new Store<RivalsLFGStore>(initialState);

export const rivalsStoreKeys = {
  all: ["groups", "profile"] as const,
  profile: (id: string) => ["profile", id] as const,
};

export const rivalsStoreActions = {
  setGroups(groups: Group[]) {
    rivalslfgStore.setState((prev) => ({ ...prev, groups }));
  },
  setProfile(profile: Profile) {
    rivalslfgStore.setState((prev) => ({ ...prev, profile }));
  },
  clearStore() {
    rivalslfgStore.setState(() => initialState);
    queryClient.removeQueries({ queryKey: rivalsStoreKeys.all });
  },
};
