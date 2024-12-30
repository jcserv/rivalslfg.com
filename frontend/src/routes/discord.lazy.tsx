import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/discord")({
  component: DiscordPage,
});

function DiscordPage() {
  return (
    <section className="min-h-[1/2] p-2 flex items-center justify-center">
      <div className="w-full flex flex-col items-center text-center">
        <div className="w-1/2">
          We&apos;re working on implementing a Discord bot for Rivals LFG,
          catered towards large Discord servers. If you have any feedback,
          please open a Github issue and we would love to learn more about what
          features you&apos;d find useful!
        </div>
      </div>
    </section>
  );
}
