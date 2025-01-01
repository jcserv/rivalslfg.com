import { useEffect, useMemo, useState } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";

import {
  GroupDisplay,
  GroupControls,
  Chat,
  AccessGroupDialog,
} from "@/components";
import { Button } from "@/components/ui";
import {
  useGroup,
  useIsAuthed,
  useJoinGroup,
  useProfile,
  useToast,
} from "@/hooks";
import { getPlayerFromProfile, Group, Profile, StatusCodes } from "@/types";

export const Route = createLazyFileRoute("/groups/$groupId")({
  component: GroupPage,
});

function GroupPage() {
  const { groupId } = Route.useParams();
  const isAuthed = useIsAuthed(groupId);
  const [g, isLoading, error] = useGroup(groupId);
  const [profile] = useProfile();
  const { toast } = useToast();

  const joinGroup = useJoinGroup();

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

  async function onJoin(p: Profile, passcode: string = "") {
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
  }

  function onLeave() {
    console.log("i'm leavin here D:");
  }

  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <div className="grid grid-cols-12 gap-4">
          <div className="col-span-8">
            <AccessGroupDialog open={showAccessDialog} onJoin={onJoin} />
            {!isLoading && (
              <GroupDisplay
                group={group}
                canUserAccessGroup={canUserAccessGroup}
                isOwner={isOwner}
              />
            )}
            <div className="flex flex-row justify-center mt-4">
              {canUserAccessGroup && (
                <div className="space-x-2">
                  {isPlayerInGroup && (
                    <Button variant="destructive" onClick={onLeave}>
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
              i love you
            </div>
          </div>
          {!isLoading && !error && (
            <div className="col-span-4 space-y-4">
              {isOwner && (
                <GroupControls
                  isGroupOpen={group?.open ?? false}
                  canUserAccessGroup={canUserAccessGroup}
                />
              )}
              <Chat canUserAccessGroup={canUserAccessGroup} />
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
