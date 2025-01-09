import { useMutation, useQuery } from "@tanstack/react-query";

import {
  createGroup,
  fetchGroup,
  fetchGroups,
  HTTPError,
  rivalsStoreKeys,
  StatusCodes,
} from "@/api";
import { Group, Profile } from "@/types";

import { usePagination } from "./paginate";
import { useProfile } from "./profile";

export function useGroups() {
  const [profile] = useProfile();

  return usePagination({
    queryKey: rivalsStoreKeys.groups,
    queryFn: async ({ paginateBy, filterBy }) => {
      const requirementsFilter = filterBy?.find(
        (f) => f.field === "areRequirementsMet",
      );

      return await fetchGroups({
        paginateBy,
        filterBy,
        playerRequirements: requirementsFilter ? profile : undefined,
      });
    },
    initialState: { pageSize: 10 },
  });
}

export function useGroup(
  id: string,
): [Group | undefined, boolean, Error | null] {
  const query = useQuery({
    queryKey: rivalsStoreKeys.group(id),
    queryFn: () => fetchGroup(id),
    staleTime: 10000,
    retry: (failureCount, error) => {
      if (
        error instanceof HTTPError &&
        error.statusCode === StatusCodes.Forbidden
      ) {
        return false; // Don't retry on 403
      }
      return failureCount < 3;
    },
  });

  return [query.data, query.isLoading, query.error];
}

type createGroupArgs = {
  profile: Profile;
};

export function useCreateGroup() {
  const { mutateAsync } = useMutation({
    mutationFn: (input: createGroupArgs) => {
      return createGroup(input.profile);
    },
  });
  return mutateAsync;
}
