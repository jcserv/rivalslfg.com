import { createLazyFileRoute } from "@tanstack/react-router";

import { ProfileForm } from "@/components";
import { useProfile, useLocalStorage } from "@/hooks";
import { FOURTEEN_DAYS_FROM_TODAY, Profile } from "@/types";

export const Route = createLazyFileRoute("/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const [profileId, setProfileId] = useLocalStorage(
    "profileId",
    "",
    FOURTEEN_DAYS_FROM_TODAY,
  );
  const [profile, isLoading] = useProfile(profileId);

  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        <div className="w-1/2">
          {!isLoading && (
            <ProfileForm
              initialValues={profile as Profile}
              setProfileId={setProfileId}
            />
          )}
        </div>
      </div>
    </section>
  );
}
