import { useMutation, useQuery } from "@tanstack/react-query";

import { fetchGroup, fetchGroups, rivalsStoreKeys, upsertGroup } from "@/api";
import { Group, Profile } from "@/types";

import { usePagination } from "./paginate";

export function useGroups() {
  return usePagination({
    queryKey: rivalsStoreKeys.groups,
    queryFn: ({ limit, offset }) => fetchGroups({ limit, offset }),
    initialState: { pageSize: 10 },
  });
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
