import { useState } from "react";

import {
  createFileRoute,
  Link,
  SearchSchemaInput,
} from "@tanstack/react-router";

import { ErrorBanner, FindGroupDialog, GroupTable } from "@/components";
import { Button } from "@/components/ui";
import { useProfile } from "@/hooks";

type BrowsePageSearchParams = {
  queue?: boolean;
};

export const Route = createFileRoute("/browse")({
  component: BrowsePage,
  validateSearch: (
    search: { queue?: boolean } & SearchSchemaInput,
  ): BrowsePageSearchParams => {
    return {
      ...(search.queue !== undefined && { queue: search.queue }),
    };
  },
});

function BrowsePage() {
  const searchParams = Route.useSearch();
  const [show, setShow] = useState(searchParams.queue ?? false);

  const [profile] = useProfile();
  const isProfileEmpty = !profile || Object.keys(profile).length === 0;

  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        {isProfileEmpty && (
          <ErrorBanner
            message="You must have your profile configured to join groups."
            className="w-full"
          >
            Click{" "}
            <Link to="/profile" className="hover:underline">
              here
            </Link>{" "}
            to configure your profile.
          </ErrorBanner>
        )}
        <div className="w-3/4">
          <div className="h-full flex-1 flex-col space-y-8 p-8">
            <div className="flex items-center justify-between space-y-2">
              <h2 className="text-2xl font-bold tracking-tight text-left">
                Groups
              </h2>
              <Button variant="success" onClick={() => setShow(true)}>
                Find Group
              </Button>
              <FindGroupDialog open={show} onClose={() => setShow(false)} />
            </div>
            <GroupTable profile={profile} isProfileEmpty={isProfileEmpty} />
          </div>
        </div>
      </div>
    </section>
  );
}
