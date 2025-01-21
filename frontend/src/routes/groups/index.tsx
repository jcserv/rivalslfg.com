import { useMemo, useState } from "react";

import {
  createFileRoute,
  Link,
  SearchSchemaInput,
  useRouter,
} from "@tanstack/react-router";
import { ColumnFiltersState, OnChangeFn } from "@tanstack/react-table";

import { ErrorBanner, FindGroupDialog, GroupTable } from "@/components";
import { Button } from "@/components/ui";
import { useProfile } from "@/hooks";

type GroupsPageSearchParams = {
  queue?: boolean;
  visibility?: string;
  region?: string;
  gamemode?: string;
  requirementsMet?: boolean;
};

export const Route = createFileRoute("/groups/")({
  component: GroupsPage,
  validateSearch: (
    search: GroupsPageSearchParams & SearchSchemaInput,
  ): GroupsPageSearchParams => {
    return {
      ...(search.queue !== undefined && { queue: search.queue }),
      ...(search.visibility !== undefined && { visibility: search.visibility }),
      ...(search.region !== undefined && { region: search.region }),
      ...(search.gamemode !== undefined && { gamemode: search.gamemode }),
      ...(search.requirementsMet !== undefined && {
        requirementsMet: search.requirementsMet,
      }),
    };
  },
});

function GroupsPage() {
  const { queue, visibility, region, gamemode, requirementsMet } =
    Route.useSearch();
  const router = useRouter();

  const [show, setShow] = useState(queue ?? false);

  const [profile] = useProfile();
  const isProfileEmpty = !profile || Object.keys(profile).length === 0;

  const filters: ColumnFiltersState = useMemo(() => {
    const initialFilters: ColumnFiltersState = [];
    if (visibility) {
      initialFilters.push({
        id: "open",
        value: [visibility],
      });
    }
    if (region) {
      initialFilters.push({
        id: "region",
        value: [region],
      });
    }
    if (gamemode) {
      initialFilters.push({
        id: "gamemode",
        value: [gamemode],
      });
    }
    if (requirementsMet !== undefined) {
      initialFilters.push({
        id: "areRequirementsMet",
        value: [requirementsMet.toString()],
      });
    }
    return initialFilters;
  }, [visibility, region, gamemode, requirementsMet]);

  const handleColumnFiltersChange: OnChangeFn<ColumnFiltersState> = (
    updaterOrValue,
  ) => {
    const initialFilters =
      typeof updaterOrValue === "function"
        ? updaterOrValue(filters)
        : updaterOrValue;

    const search: GroupsPageSearchParams = {
      ...(queue && { queue }),
    };

    initialFilters.forEach((filter) => {
      const filterVal = filter.value as string[];
      switch (filter.id) {
        case "open":
          if (filterVal) search.visibility = filterVal[0] as string;
          break;
        case "region":
          if (filterVal.length) search.region = filterVal[0] as string;
          break;
        case "gamemode":
          if (filterVal.length) search.gamemode = filterVal[0] as string;
          break;
        case "areRequirementsMet":
          if (filterVal.length)
            search.requirementsMet = filterVal[0] === "true";
          break;
      }
    });

    router.navigate({
      to: "/groups",
      search,
      replace: true,
    });
  };

  const handleClearFilters = () => {
    router.navigate({
      to: "/groups",
      search: queue ? { queue } : {},
      replace: true,
    });
  };

  return (
    <section className="p-2 md:p-4">
      <div className="w-full flex flex-col items-center">
        {isProfileEmpty && (
          <ErrorBanner
            message="You must have your profile configured to join groups."
            className="w-full"
          >
            Click{" "}
            <Link to="/profile" className="hover:underline">
              here
            </Link>{" "}
            to configure your profile.
          </ErrorBanner>
        )}
        <div className="w-3/4">
          <div className="h-full flex-1 flex-col space-y-8 p-8">
            <div className="flex items-center justify-between space-y-2">
              <h2 className="text-2xl font-bold tracking-tight text-left">
                Groups
              </h2>
              <Button variant="success" onClick={() => setShow(true)}>
                Find Group
              </Button>
              <FindGroupDialog open={show} onClose={() => setShow(false)} />
            </div>
            <GroupTable
              profile={profile}
              isProfileEmpty={isProfileEmpty}
              filters={filters}
              handleColumnFiltersChange={handleColumnFiltersChange}
              handleClearFilters={handleClearFilters}
            />
          </div>
        </div>
      </div>
    </section>
  );
}
