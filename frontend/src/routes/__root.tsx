import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRootRoute, Outlet } from "@tanstack/react-router";

import { RivalsLFGClient } from "@/api";
import { ErrorBanner, Footer, Header } from "@/components";
import { Toaster } from "@/components/ui";

export const queryClient = new QueryClient();
export const rivalslfgAPIClient = new RivalsLFGClient(
  import.meta.env.VITE_API_URL,
);

export const Route = createRootRoute({
  component: () => (
    <>
      <QueryClientProvider client={queryClient}>
        <ReactQueryDevtools initialIsOpen={false} />
        <Header />
        <hr />
        <Outlet />
        <Footer />
        <Toaster />
      </QueryClientProvider>
    </>
  ),
  errorComponent: ({ error }) => (
    <>
      <Header />
      <hr />
      <div className="flex w-full h-[80vh] justify-center align-middle">
        <ErrorBanner
          message="An unexpected error occurred. Please open a Github issue with the below error details."
          error={error.message}
          className="h-1/3"
        />
      </div>
      <Footer />
    </>
  ),
});
