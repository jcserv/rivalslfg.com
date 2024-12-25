import { Footer } from "@/components/Footer";
import { Header } from "@/components/Header";
import { createRootRoute, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
  component: () => (
    <>
      <Header />
      <hr />
      <Outlet />
      <Footer />
    </>
  ),
});
