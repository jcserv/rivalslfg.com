import { useMutation, useQuery } from "@tanstack/react-query";

import {
  createGroup,
  deleteGroup,
  fetchGroup,
  fetchGroups,
  HTTPError,
  rivalsStoreKeys,
  StatusCodes,
} from "@/api";
import { queryClient } from "@/routes/__root";
import { Group, Player, Profile } from "@/types";

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

export function addPlayerToGroup(groupId: string, player: Player) {
  queryClient.setQueryData<Group>(
    rivalsStoreKeys.group(groupId),
    (oldGroup) => {
      if (!oldGroup) return;
      return {
        ...oldGroup,
        players: [...oldGroup.players, player],
      };
    },
  );
}

export function removePlayerFromGroup(
  groupId: string,
  playerId: number,
  newLeaderId: number,
) {
  queryClient.setQueryData<Group>(
    rivalsStoreKeys.group(groupId),
    (oldGroup) => {
      if (!oldGroup) return;
      const newPlayers = oldGroup.players.reduce(
        (acc, player) => {
          if (player.id === playerId) {
            return acc;
          }

          return [
            ...acc,
            {
              ...player,
              leader: player.id === newLeaderId,
            },
          ];
        },
        [] as typeof oldGroup.players,
      );

      return {
        ...oldGroup,
        ownerId: newLeaderId,
        players: newPlayers,
      };
    },
  );
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

export function useDeleteGroup() {
  const { mutateAsync } = useMutation({
    mutationFn: (id: string) => {
      return deleteGroup(id);
    },
  });
  return mutateAsync;
}