import { createLazyFileRoute, Link } from "@tanstack/react-router";

import { ErrorBanner, GroupTable } from "@/components";
import { useProfile } from "@/hooks";

export const Route = createLazyFileRoute("/browse")({
  component: BrowsePage,
});

function BrowsePage() {
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
          <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
            <div className="flex items-center justify-between space-y-2">
              <h2 className="text-2xl font-bold tracking-tight text-left">
                Groups
              </h2>
            </div>
            <GroupTable profile={profile} isProfileEmpty={isProfileEmpty} />
          </div>
        </div>
      </div>
    </section>
  );
}
