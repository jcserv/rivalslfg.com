import { Table } from "@tanstack/react-table";
import { X } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

import {
  areRequirementsMet,
  gamemodes,
  open,
  platforms,
  regions,
} from "./data";
import { DataTableFacetedFilter } from "./data-table-faceted-filter";

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  rightAdornment?: React.ReactNode;
}

export function DataTableToolbar<TData>({
  table,
  rightAdornment,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;
  return (
    <div className="flex items-center justify-between overflow-x-auto">
      <div className="flex flex-1 items-center space-x-2">
        <Input
          placeholder="Search groups..."
          value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("name")?.setFilterValue(event.target.value)
          }
          className="h-8 w-[150px] lg:w-[175px]"
        />
        {table.getColumn("open") && (
          <DataTableFacetedFilter
            column={table.getColumn("open")}
            title="Visibility"
            options={open}
          />
        )}
        {table.getColumn("region") && (
          <DataTableFacetedFilter
            column={table.getColumn("region")}
            title="Region"
            options={regions}
          />
        )}
        {table.getColumn("gamemode") && (
          <DataTableFacetedFilter
            column={table.getColumn("gamemode")}
            title="Gamemode"
            options={gamemodes}
          />
        )}
        {/* {table.getColumn("platform") && (
          <DataTableFacetedFilter
            column={table.getColumn("platform")}
            title="Platform"
            options={platforms}
          />
        )} */}
        {table.getColumn("areRequirementsMet") && (
          <DataTableFacetedFilter
            column={table.getColumn("areRequirementsMet")}
            title="Requirements Met"
            options={areRequirementsMet}
          />
        )}
        {isFiltered && (
          <Button
            variant="ghost"
            onClick={() => table.resetColumnFilters()}
            className="h-8 px-2 lg:px-3"
          >
            Clear
            <X />
          </Button>
        )}
      </div>
      {rightAdornment}
    </div>
  );
}
