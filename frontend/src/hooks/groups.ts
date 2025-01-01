import { useMutation, useQuery } from "@tanstack/react-query";

import {
  rivalsStoreKeys,
  fetchGroups,
  fetchGroup,
  joinGroup,
  upsertGroup,
} from "@/api";
import { Group, Profile } from "@/types";

export function useGroups(): [Group[] | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.groups,
    queryFn: () => fetchGroups(),
    staleTime: 0,
  });

  return [query.data, query.isLoading, query.error];
}

export function useGroup(
  id: string,
): [Group | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.group(id),
    queryFn: () => fetchGroup(id),
    staleTime: 0,
  });

  return [query.data, query.isLoading, query.error];
}

type UpsertGroupArgs = {
  profile: Profile;
  id: string;
};

export function useUpsertGroup() {
  const { mutateAsync } = useMutation({
    mutationFn: (input: UpsertGroupArgs) => {
      return upsertGroup(input.profile, input.id);
    },
  });
  return mutateAsync;
}

type JoinGroupArgs = {
  groupId: string;
  player: Profile;
  passcode: string;
};

export function useJoinGroup() {
  const { mutateAsync } = useMutation({
    mutationFn: (input: JoinGroupArgs) => {
      return joinGroup(input.groupId, input.player, input.passcode);
    },
  });
  return mutateAsync;
}
