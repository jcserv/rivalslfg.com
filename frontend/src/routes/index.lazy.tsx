import { InfoBanner } from "@/components/Banner";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <section className="p-2 md:p-4">
      <div className="min-h-[80vh] w-full flex flex-col items-center">
        <InfoBanner><p>This is a test page</p></InfoBanner>
      </div>
    </section>
  );
}
