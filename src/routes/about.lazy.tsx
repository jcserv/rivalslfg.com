import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/about")({
  component: About,
});

function About() {
  return <section className="p-2">I&apos;m a dog with a blog.</section>;
}
