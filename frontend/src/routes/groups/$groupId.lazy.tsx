import { createLazyFileRoute } from "@tanstack/react-router";

import groups from "@/assets/groups.json";
import { GroupDisplay } from "@/components/GroupDisplay";
import { GroupControls } from "@/components/GroupControls";
import { Group as GroupType } from "@/types";
import { Chat } from "@/components/Chat";
import { Button } from "@/components/ui";

export const Route = createLazyFileRoute("/groups/$groupId")({
  component: Group,
});

function Group() {
  const { groupId } = Route.useParams();
  const group = groups.find((group) => group.id === groupId);

  function onLeave() {
    // TODO: This should also be logged in the chat
    console.log("i'm leavin here D:");
  }

  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <div className="grid grid-cols-12 gap-4">
          <div className="col-span-8">
            <GroupDisplay group={group as GroupType} />
            <div className="flex flex-row justify-center mt-4">
              <Button variant="destructive" onClick={onLeave}>
                Leave
              </Button>
            </div>
          </div>
          <div className="col-span-4 space-y-4">
            <GroupControls isGroupOpen={group?.open || false} />
            <Chat />
          </div>
        </div>
      </div>
    </section>
  );
}
