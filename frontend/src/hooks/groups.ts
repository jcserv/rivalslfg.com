import { useQuery } from "@tanstack/react-query";

import { rivalsStoreKeys, fetchGroups, fetchGroup } from "@/api";
import { Group, ONE_MINUTE_IN_MILLISECONDS } from "@/types";

export function useGroups(): [Group[] | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.all,
    queryFn: () => fetchGroups(),
    staleTime: ONE_MINUTE_IN_MILLISECONDS * 5,
  });

  return [query.data, query.isLoading, query.error];
}

export function useGroup(
  id: string,
): [Group | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.group(id),
    queryFn: () => fetchGroup(id),
    staleTime: ONE_MINUTE_IN_MILLISECONDS * 5,
  });

  return [query.data, query.isLoading, query.error];
}
