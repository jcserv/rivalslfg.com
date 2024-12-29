import { ProfileForm } from "@/components";
import { useLocalStorage } from "@/hooks/localStorage";
import { Profile } from "@/types";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const fourteenDaysFromToday = new Date(
    new Date().setDate(new Date().getDate() + 14),
  );
  const [profile, setProfile] = useLocalStorage(
    "profile",
    {},
    fourteenDaysFromToday,
  );

  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        <div className="w-1/2">
          <ProfileForm
            initialValues={profile as Profile}
            setProfile={setProfile}
          />
        </div>
      </div>
    </section>
  );
}
