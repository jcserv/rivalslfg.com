import { InfoBanner, ProfileForm } from "@/components";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui";
import { useProfile } from "@/hooks";
import { useLocalStorage } from "@/hooks/localStorage";
import { FOURTEEN_DAYS_FROM_TODAY, Profile } from "@/types";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/")({
  component: Index,
});

function Index() {
  const [profileId, setProfileId] = useLocalStorage(
    "profileId",
    "",
    FOURTEEN_DAYS_FROM_TODAY,
  );
  const [profile, isLoading] = useProfile(profileId);

  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        <InfoBanner>
          <p>
            Rivals LFG is a platform for Marvel Rivals players to find groups to
            play with - based on your rank, region, platform, etc. It&apos;ll
            also suggest team-ups for your group to play.
          </p>
        </InfoBanner>
        <Tabs defaultValue="find" className="w-1/2">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="find">Find Group</TabsTrigger>
            <TabsTrigger value="create">Create Group</TabsTrigger>
          </TabsList>
          <TabsContent value="find">
            {!isLoading && (
              <ProfileForm
                initialValues={profile as Profile}
                setProfileId={setProfileId}
              />
            )}
          </TabsContent>
          <TabsContent value="create">
            {!isLoading && (
              <ProfileForm
                isGroup
                initialValues={profile as Profile}
                setProfileId={setProfileId}
              />
            )}
          </TabsContent>
        </Tabs>
      </div>
    </section>
  );
}
