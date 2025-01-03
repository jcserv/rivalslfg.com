import { createLazyFileRoute } from "@tanstack/react-router";

import { InfoBanner, ProfileForm } from "@/components";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui";
import { useProfile } from "@/hooks";
import { Profile } from "@/types";

export const Route = createLazyFileRoute("/")({
  component: Index,
});

function Index() {
  const [profile, setProfile] = useProfile();

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
            <ProfileForm
              profileFormType="find"
              profile={profile as Profile}
              setProfile={setProfile}
            />
          </TabsContent>
          <TabsContent value="create">
            <ProfileForm
              profileFormType="create"
              profile={profile as Profile}
              setProfile={setProfile}
            />
          </TabsContent>
        </Tabs>
      </div>
    </section>
  );
}
