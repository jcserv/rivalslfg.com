import { useCallback, useEffect, useMemo, useState } from "react";

import {
  createFileRoute,
  notFound,
  SearchSchemaInput,
  useRouter,
} from "@tanstack/react-router";

import { fetchGroup, joinGroup } from "@/api";
import { HTTPError, StatusCodes } from "@/api/types";
import {
  AccessGroupDialog,
  ChatBox,
  GroupControls,
  GroupDisplay,
  TooltipItem,
} from "@/components";
import { Button } from "@/components/ui";
import {
  getStorageValue,
  setStorageValue,
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
  beforeLoad: async ({ params, search }) => {
    const { groupId } = params;
    const { join, passcode } = search;
    try {
      await fetchGroup(groupId.toUpperCase());
      return;
    } catch (error) {
      if (!(error instanceof HTTPError)) {
        return;
      }
      if (error.statusCode === StatusCodes.NotFound) {
        throw notFound();
      }
      if (error.statusCode === StatusCodes.Forbidden) {
        if (!join) return;
        const profile: Profile = getStorageValue(
          "profile",
          {} as Profile,
        ).value;
        const { playerId } = await joinGroup(
          groupId.toUpperCase(),
          profile,
          passcode ?? "",
        );
        setStorageValue("profile", {
          ...profile,
          id: playerId,
        });
      }
    }
  },
});

function GroupPage() {
  const { groupId } = Route.useParams();
  const { toast } = useToast();
  const router = useRouter();

  const isAuthed = useIsAuthed(groupId);
  const [g, isLoading, error] = useGroup(groupId);
  const [profile, setProfile, isProfileConfigured] = useProfile();

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
        const { playerId } = await joinGroup({
          groupId,
          player: p,
          passcode,
        });

        setProfile({
          ...profile,
          id: playerId,
        });

        if (!isPlayerInGroup) {
          setGroup({
            ...group,
            players: [...group.players, getPlayerFromProfile(p, playerId)],
          });
        }

        setShowAccessDialog(false);
        setCanUserAccessGroup(true);

        toast({
          title: "Joined group",
          variant: "success",
        });
      } catch (error) {
        if (!(error instanceof HTTPError)) {
          toast({
            title: "Access denied",
            description: "Please try again.",
            variant: "destructive",
          });
          return;
        }
        toast({
          title: error.statusText,
          description: error.message,
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

  async function onRemove(playerId: number) {
    if (!group) return;
    const isPlayerLeavingGroup = playerId === profile.id;
    try {
      const status = await removePlayer({
        groupId,
        playerId,
      });

      toast({
        title: isPlayerLeavingGroup
          ? "Left group"
          : "Removed player from group",
        variant: "success",
      });

      // Group was deleted
      if (status === StatusCodes.NoContent) {
        router.navigate({ to: "/groups" });
        return;
      }

      const newGroup = {
        ...group,
        players: group.players
          ? group.players.filter((p) => p.id !== playerId)
          : [],
      };
      setGroup(newGroup);
    } catch {
      toast({
        title: isPlayerLeavingGroup
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
                passcode={group?.passcode ?? ""}
                onRemove={onRemove}
              />
            )}
            <div className="flex flex-row justify-center mt-4">
              {canUserAccessGroup && (
                <div className="space-x-2">
                  {isPlayerInGroup && (
                    <Button
                      variant="destructive"
                      onClick={() => onRemove(profile.id ?? 0)}
                    >
                      Leave
                    </Button>
                  )}
                  {!isPlayerInGroup && (
                    <TooltipItem
                      content="You must configure your profile before joining a group"
                      enabled={!isProfileConfigured}
                    >
                      <Button
                        variant="success"
                        disabled={!isProfileConfigured}
                        onClick={(e) => {
                          e.preventDefault();
                          onJoin(profile);
                        }}
                      >
                        Join
                      </Button>
                    </TooltipItem>
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
                  passcode={group?.passcode ?? ""}
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
