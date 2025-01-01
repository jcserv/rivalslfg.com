import { createLazyFileRoute, Link } from "@tanstack/react-router";

import { ErrorBanner, GroupTable } from "@/components";
import { useGroups, useProfile } from "@/hooks";

export const Route = createLazyFileRoute("/browse")({
  component: BrowsePage,
});

function BrowsePage() {
  const [profile] = useProfile();
  const [groups, isLoadingGroups] = useGroups();

  const isProfileEmpty = !profile || Object.keys(profile).length === 0;
  const isLoading = isLoadingGroups;

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
          <GroupTable
            groups={!isLoading && groups ? groups : []}
            profile={profile}
            isProfileEmpty={isProfileEmpty}
          />
        </div>
      </div>
    </section>
  );
}
