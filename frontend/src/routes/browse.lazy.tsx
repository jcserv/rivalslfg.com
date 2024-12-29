import { GroupTable } from "@/components/GroupTable";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/browse")({
  component: BrowsePage,
});

function BrowsePage() {
  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        <div className="w-3/4">
          <GroupTable />
        </div>
      </div>
    </section>
  );
}
