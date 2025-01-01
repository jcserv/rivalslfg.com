import { useStore } from "@tanstack/react-store";

import { rivalslfgStore } from "@/api";

export const useIsAuthed = (groupId: string): boolean => {
  return useStore(
    rivalslfgStore,
    (state) => state["authedGroups"].get(groupId) ?? false,
  );
};
