import { useCallback, useEffect, useMemo, useState } from "react";

import { createFileRoute, SearchSchemaInput } from "@tanstack/react-router";

import { StatusCodes } from "@/api/types";
import {
  AccessGroupDialog,
  ChatBox,
  GroupControls,
  GroupDisplay,
} from "@/components";
import { Button } from "@/components/ui";
import {
  useGroup,
  useIsAuthed,
  useJoinGroup,
  useProfile,
  useRemovePlayer,
  useToast,
} from "@/hooks";
import { getPlayerFromProfile, Group, Profile } from "@/types";

type GroupPageSearchParams = {
  join?: boolean;
  passcode?: string;
};

export const Route = createFileRoute("/groups/$groupId")({
  component: GroupPage,
  validateSearch: (
    search: { join?: boolean; passcode?: string } & SearchSchemaInput,
  ): GroupPageSearchParams => {
    return {
      ...(search.join !== undefined && { join: search.join }),
      ...(search.passcode !== undefined && { passcode: search.passcode }),
    };
  },
});

function GroupPage() {
  const passcode = "abcd";

  const { groupId } = Route.useParams();
  const searchParams = Route.useSearch();

  const isAuthed = useIsAuthed(groupId);
  const [g, isLoading, error] = useGroup(groupId);
  const [profile] = useProfile();
  const { toast } = useToast();

  const joinGroup = useJoinGroup();
  const removePlayer = useRemovePlayer();

  const [group, setGroup] = useState<Group | undefined>(g);
  const [showAccessDialog, setShowAccessDialog] = useState(false);
  const [canUserAccessGroup, setCanUserAccessGroup] = useState<boolean | null>(
    null,
  );

  const { isPlayerInGroup, isOwner } = useMemo(() => {
    return {
      isPlayerInGroup: group?.players.some((p) => p.name === profile.name),
      isOwner: group?.owner === profile.name,
    };
  }, [group, profile]);

  useEffect(() => {
    if (g) {
      setGroup(g);
      const hasAccess = g.open || isAuthed;
      setCanUserAccessGroup(hasAccess);
      setShowAccessDialog(!isLoading && g && !hasAccess);
    }
  }, [g, isLoading, isAuthed]);

  const onJoin = useCallback(
    async (p: Profile, passcode: string = "") => {
      if (!group) return;
      try {
        const status = await joinGroup({
          groupId,
          player: p,
          passcode,
        });
        if (status !== StatusCodes.OK) {
          throw new Error(`${status}`);
        }
        if (!isPlayerInGroup) {
          setGroup({
            ...group,
            players: [...group.players, getPlayerFromProfile(p)],
          });
        }

        setShowAccessDialog(false);
        setCanUserAccessGroup(true);

        toast({
          title: "Joined group",
          variant: "success",
        });
      } catch {
        toast({
          title: "Access denied",
          description: "Please try again.",
          variant: "destructive",
        });
      }
    },
    [
      group,
      groupId,
      isPlayerInGroup,
      setGroup,
      setShowAccessDialog,
      setCanUserAccessGroup,
    ],
  );

  useEffect(() => {
    if (searchParams.join) {
      const joinGroupAsync = async () => {
        await onJoin(profile, searchParams.passcode ?? "");
      };
      joinGroupAsync();
    }
  }, [searchParams.join, profile, searchParams.passcode, onJoin]);

  async function onRemove(id: number, playerToRemove: string) {
    if (!group) return;
    try {
      const status = await removePlayer({
        groupId,
        playerId: 1, // TODO: This should be generated server-side
        requesterName: profile.name,
        playerName: playerToRemove,
      });
      if (status !== StatusCodes.OK) {
        throw new Error(`${status}`);
      }

      if (isPlayerInGroup) {
        const updatedPlayers = group.players.filter(
          (p) => p.name !== playerToRemove,
        );

        let newGroup = {
          ...group,
          players: updatedPlayers,
        };

        if (group.owner === playerToRemove) {
          const newOwner = updatedPlayers[0];

          newGroup = {
            ...newGroup,
            owner: newOwner.name,
            name: `${newOwner.name}'s group`,
            players: updatedPlayers.map((p) =>
              p.name === newOwner.name ? { ...p, leader: true } : p,
            ),
          };
        }

        setGroup(newGroup);
      }

      toast({
        title:
          playerToRemove === profile.name
            ? "Left group"
            : "Removed player from group",
        variant: "success",
      });
    } catch {
      toast({
        title:
          playerToRemove === profile.name
            ? "Unable to leave group"
            : "Unable to remove player from group",
        description: "Please try again.",
        variant: "destructive",
      });
    }
  }

  return (
    <section className="p-2 md:p-4">
      <div className="h-full w-full flex flex-col items-center">
        <div className="grid grid-cols-1 md:grid-cols-12 gap-4 w-full max-w-[1200px]">
          <div className="col-span-1 md:col-span-8">
            <AccessGroupDialog open={showAccessDialog} onJoin={onJoin} />
            {!isLoading && (
              <GroupDisplay
                group={group}
                canUserAccessGroup={canUserAccessGroup}
                isOwner={isOwner}
                passcode={passcode}
                onRemove={onRemove}
              />
            )}
            <div className="flex flex-row justify-center mt-4">
              {canUserAccessGroup && (
                <div className="space-x-2">
                  {isPlayerInGroup && (
                    <Button
                      variant="destructive"
                      onClick={() => onRemove(profile.id, profile.name)}
                    >
                      Leave
                    </Button>
                  )}
                  {!isPlayerInGroup && (
                    <Button
                      variant="success"
                      onClick={(e) => {
                        e.preventDefault();
                        onJoin(profile);
                      }}
                    >
                      Join
                    </Button>
                  )}
                </div>
              )}
            </div>
          </div>
          {!isLoading && !error && (
            <div className="col-span-1 md:col-span-4 space-y-4 sm:w-1/2  md:w-full mx-auto">
              {isOwner && (
                <GroupControls
                  isGroupOpen={group?.open ?? false}
                  canUserAccessGroup={canUserAccessGroup}
                  passcode={passcode}
                />
              )}
              <ChatBox
                canUserAccessGroup={canUserAccessGroup}
                isPlayerInGroup={isPlayerInGroup}
              />
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
