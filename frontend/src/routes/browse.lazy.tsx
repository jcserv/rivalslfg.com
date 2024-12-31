import { createLazyFileRoute, Link } from "@tanstack/react-router";

import { useLocalStorage, useGroups, useProfile } from "@/hooks";
import { ErrorBanner, GroupTable } from "@/components";
import { Group, FOURTEEN_DAYS_FROM_TODAY } from "@/types";

export const Route = createLazyFileRoute("/browse")({
  component: BrowsePage,
});

function BrowsePage() {
  const [profileId] = useLocalStorage(
    "profileId",
    "-1",
    FOURTEEN_DAYS_FROM_TODAY,
  );
  const [profile, isLoadingProfile] = useProfile(profileId);
  const [groups, isLoadingGroups] = useGroups();

  const isProfileEmpty = !profile || Object.keys(profile).length === 0;
  const isLoading = isLoadingProfile || isLoadingGroups;

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
          {!isLoading && groups && (
            <GroupTable
              groups={groups as Group[]}
              profile={profile}
              isProfileEmpty={isProfileEmpty}
            />
          )}
        </div>
      </div>
    </section>
  );
}
