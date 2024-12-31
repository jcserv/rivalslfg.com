import { useQuery } from "@tanstack/react-query";

import { rivalsStoreKeys, fetchProfile } from "@/api";
import { ONE_MINUTE_IN_MILLISECONDS, Profile } from "@/types";

export function useProfile(
  id: string,
): [Profile | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.profile(id),
    queryFn: () => fetchProfile(id),
    staleTime: ONE_MINUTE_IN_MILLISECONDS * 60,
  });

  return [query.data, query.isLoading, query.error];
}
