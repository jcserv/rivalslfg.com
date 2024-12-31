import { useQuery } from "@tanstack/react-query";

import { rivalsStoreKeys, fetchGroups } from "@/api";
import { ONE_MINUTE_IN_MILLISECONDS } from "@/types";

export function useGroups() {
  const query = useQuery({
    queryKey: rivalsStoreKeys.all,
    queryFn: () => fetchGroups(),
    staleTime: ONE_MINUTE_IN_MILLISECONDS * 5,
  });

  return [query.data, query.isLoading, query.error];
}