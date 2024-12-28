import { ProfileForm } from "@/components";
import { InfoBanner } from "@/components/Banner";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/profile")({
  component: Profile,
});

function Profile() {
  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <InfoBanner>
          <p>This is a test page</p>
        </InfoBanner>
        <ProfileForm />
      </div>
    </section>
  );
}
