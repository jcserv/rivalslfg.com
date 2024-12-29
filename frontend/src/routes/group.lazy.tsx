import { createLazyFileRoute } from "@tanstack/react-router";

import mockGroup from "@/assets/mockGroup.json";
import { GroupDisplay } from "@/components/GroupDisplay";
import { Group as GroupType } from "@/types";

export const Route = createLazyFileRoute("/group")({
  component: Group,
});

function Group() {
  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <GroupDisplay group={mockGroup as GroupType} />
      </div>
    </section>
  );
}
