import { createLazyFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";

import { Group as GroupType } from "@/types";
import { Button } from "@/components/ui";
import { useGroup } from "@/hooks";
import {
  GroupDisplay,
  GroupControls,
  Chat,
  AccessGroupDialog,
} from "@/components";

export const Route = createLazyFileRoute("/groups/$groupId")({
  component: Group,
});

function Group() {
  const { groupId } = Route.useParams();
  const [group, isLoading, error] = useGroup(groupId);
  const isGroupOpen = group?.open || false;

  const [canUserAccessGroup, setCanUserAccessGroup] = useState(false);
  useEffect(() => {
    setCanUserAccessGroup(isGroupOpen);
  }, [isGroupOpen]);

  function onLeave() {
    // TODO: This should also be logged in the chat
    console.log("i'm leavin here D:");
  }

  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <div className="grid grid-cols-12 gap-4">
          <div className="col-span-8">
            <AccessGroupDialog
              open={!canUserAccessGroup}
              onClose={() => {
                setCanUserAccessGroup(true);
              }}
            />
            {!isLoading && (
              <GroupDisplay
                group={group as GroupType}
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
