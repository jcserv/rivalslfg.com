import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <section className="p-2 text-center">
      <div className="h-screen w-screen flex items-center justify-center">
        <h3>Welcome to my blog!</h3>
      </div>
    </section>
  );
}
