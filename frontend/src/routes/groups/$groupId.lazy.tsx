import { useEffect, useState } from "react";
import { createLazyFileRoute } from "@tanstack/react-router";

import {
  GroupDisplay,
  GroupControls,
  Chat,
  AccessGroupDialog,
} from "@/components";
import { Button } from "@/components/ui";
import { useGroup } from "@/hooks";
import { getPlayerFromProfile, Group, Profile } from "@/types";
import { rivalsStoreActions } from "@/api";

export const Route = createLazyFileRoute("/groups/$groupId")({
  component: GroupPage,
});

function GroupPage() {
  const { groupId } = Route.useParams();
  const isAuthed = rivalsStoreActions.getIsAuthed(groupId);

  // Initialize access state as null to represent "unknown" state
  const [canUserAccessGroup, setCanUserAccessGroup] = useState<boolean | null>(
    null,
  );

  const [g, isLoading, error] = useGroup(groupId);
  const [group, setGroup] = useState<Group | undefined>(g);
  const isGroupOpen = group?.open || false;

  useEffect(() => {
    // Only set access state when we have loaded the group data
    if (!isLoading) {
      setCanUserAccessGroup(isGroupOpen || isAuthed);
    }
  }, [isLoading, isGroupOpen]);

  function onJoin(p: Profile) {
    if (!group) return;
    setGroup({
      ...group,
      players: [...group.players, getPlayerFromProfile(p)],
    });
  }

  function onLeave() {
    // TODO: This should also be logged in the chat
    console.log("i'm leavin here D:");
  }

  const accessStateUnknown = isLoading || canUserAccessGroup === null;
  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <div className="grid grid-cols-12 gap-4">
          <div className="col-span-8">
            <AccessGroupDialog
              groupId={groupId}
              open={!accessStateUnknown && !canUserAccessGroup}
              onJoin={onJoin}
              onClose={() => {
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
                isGroupOpen={isGroupOpen}
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
