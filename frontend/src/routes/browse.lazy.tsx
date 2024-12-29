import { GroupTable } from "@/components/GroupTable";
import { createLazyFileRoute } from "@tanstack/react-router";

import groups from "@/assets/groups.json";
import { Group } from "@/types";

export const Route = createLazyFileRoute("/browse")({
  component: BrowsePage,
});

function BrowsePage() {
  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        <div className="w-3/4">
          <GroupTable groups={groups as Group[]} />
        </div>
      </div>
    </section>
  );
}
