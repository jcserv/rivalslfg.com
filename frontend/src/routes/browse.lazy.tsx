import { GroupTable } from "@/components/GroupTable";
import { createLazyFileRoute, Link } from "@tanstack/react-router";

import groups from "@/assets/groups.json";
import { Group } from "@/types";
import { useLocalStorage } from "@/hooks/localStorage";
import { ErrorBanner } from "@/components";

export const Route = createLazyFileRoute("/browse")({
  component: BrowsePage,
});

function BrowsePage() {
  const fourteenDaysFromToday = new Date(
    new Date().setDate(new Date().getDate() + 14),
  );
  const [profile] = useLocalStorage("profile", {}, fourteenDaysFromToday);
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
          <GroupTable
            groups={groups as Group[]}
            profile={profile}
            isProfileEmpty={isProfileEmpty}
          />
        </div>
      </div>
    </section>
  );
}
