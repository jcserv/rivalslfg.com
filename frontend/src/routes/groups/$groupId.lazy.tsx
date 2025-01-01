import { useEffect, useState } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";

import {
  GroupDisplay,
  GroupControls,
  Chat,
  AccessGroupDialog,
} from "@/components";
import { Button } from "@/components/ui";
import { useGroup, useIsAuthed } from "@/hooks";
import { getPlayerFromProfile, Group, Profile } from "@/types";

export const Route = createLazyFileRoute("/groups/$groupId")({
  component: GroupPage,
});

function GroupPage() {
  const { groupId } = Route.useParams();
  const isAuthed = useIsAuthed(groupId);
  const [g, isLoading, error] = useGroup(groupId);
  const [group, setGroup] = useState<Group | undefined>(g);
  const [showAccessDialog, setShowAccessDialog] = useState(false);
  const [canUserAccessGroup, setCanUserAccessGroup] = useState<boolean | null>(null);

  useEffect(() => {
    if (g) {
      setGroup(g);
      const hasAccess = g.open || isAuthed;
      setCanUserAccessGroup(hasAccess);
      setShowAccessDialog(!isLoading && g && !hasAccess);
    }
  }, [g, isLoading, isAuthed]);

  function onJoin(p: Profile) {
    if (!group) return;
    setGroup({
      ...group,
      players: [...group.players, getPlayerFromProfile(p)],
    });
  }

  function onLeave() {
    console.log("i'm leavin here D:");
  }

  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <div className="grid grid-cols-12 gap-4">
          <div className="col-span-8">
            <AccessGroupDialog
              groupId={groupId}
              open={showAccessDialog}
              onJoin={onJoin}
              onClose={() => {
                setShowAccessDialog(false);
                setCanUserAccessGroup(true);
              }}
            />
            {!isLoading && (
              <GroupDisplay
                group={group}
                canUserAccessGroup={canUserAccessGroup}
              />
            )}
            <div className="flex flex-row justify-center mt-4">
              {canUserAccessGroup && (
                <Button variant="destructive" onClick={onLeave}>
                  Leave
                </Button>
              )}
            </div>
          </div>
          {!isLoading && !error && (
            <div className="col-span-4 space-y-4">
              <GroupControls
                isGroupOpen={group?.open ?? false}
                canUserAccessGroup={canUserAccessGroup}
              />
              <Chat canUserAccessGroup={canUserAccessGroup} />
            </div>
          )}
        </div>
      </div>
    </section>
  );
}