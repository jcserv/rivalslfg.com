import { Header, Footer } from "@/components";
import { Toaster } from "@/components/ui";
import { createRootRoute, Outlet } from "@tanstack/react-router";

/*
Root (/): Displays the info about the application, with a button to find a lobby or create a new one.
Profile (/profile): Displays the profile form. User is redirected to this page if they have not filled it out yet.
Find Lobby (/lobby/find): Displays the lobby search form. User is redirected to this page if they select to find a lobby.
Create Lobby (/lobby/create): Displays the lobby creation form. User is redirected to this page if they select to create a lobby.
*/

export const Route = createRootRoute({
  component: () => (
    <>
      <Header />
      <hr />
      <Outlet />
      <Footer />
      <Toaster />
    </>
  ),
});
