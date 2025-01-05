import { useMutation } from "@tanstack/react-query";

import { joinGroup, removePlayer } from "@/api";
import { Profile } from "@/types";

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

type RemovePlayerArgs = {
  groupId: string;
  playerId: number;
};

export function useRemovePlayer() {
  const { mutateAsync } = useMutation({
    mutationFn: (input: RemovePlayerArgs) => {
      return removePlayer(
        input.groupId,
        input.playerId,
      );
    },
  });
  return mutateAsync;
}
