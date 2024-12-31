import { useMutation, useQuery } from "@tanstack/react-query";

import { rivalsStoreKeys, fetchProfile } from "@/api";
import { ONE_MINUTE_IN_MILLISECONDS, Profile } from "@/types";
import { rivalslfgAPIClient } from "@/routes/__root";

type CreateProfileArgs = {
  profile: Profile;
  id: number | undefined;
};

export function useCreateProfile() {
  const { mutateAsync } = useMutation({
    mutationFn: (input: CreateProfileArgs) => {
      return rivalslfgAPIClient.createPlayer(input.profile, input.id);
    },
  });
  return mutateAsync;
}

export function useProfile(
  id: string,
): [Profile | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.profile(id),
    queryFn: () => fetchProfile(id),
    staleTime: ONE_MINUTE_IN_MILLISECONDS * 60,
    retry: false,
  });

  return [query.data, query.isLoading, query.error];
}
