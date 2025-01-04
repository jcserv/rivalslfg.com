import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRootRoute, Outlet } from "@tanstack/react-router";

import { RivalsLFGClient } from "@/api/rivalslfg";
import { BackButton, ErrorBanner, Footer, Header } from "@/components";
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
  notFoundComponent: () => (
    <>
      <hr />
      <div className="flex w-full h-[80vh] justify-center align-middle">
        <ErrorBanner
          message="Page not found. Please check the URL and try again."
          className="h-[125px]"
        >
          <BackButton link className="p-2" />
        </ErrorBanner>
      </div>
    </>
  ),
  errorComponent: ({ error }) => (
    <>
      <hr />
      <div className="flex w-full h-[80vh] justify-center align-middle">
        <ErrorBanner
          message="An unexpected error occurred. Please open a Github issue with the below error details."
          error={error.message}
          className="h-1/3"
        />
      </div>
    </>
  ),
});
