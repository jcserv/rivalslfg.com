import { createLazyFileRoute } from "@tanstack/react-router";

import { ProfileForm } from "@/components";
import { useProfile } from "@/hooks";
import { Profile } from "@/types";

export const Route = createLazyFileRoute("/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const [profile, setProfile] = useProfile();

  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        <div className="w-1/2">
          <ProfileForm
            profileFormType="profile"
            profile={profile as Profile}
            setProfile={setProfile}
          />
        </div>
      </div>
    </section>
  );
}
