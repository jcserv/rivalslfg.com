import { ColumnDef, Row } from "@tanstack/react-table";

import { DataTableColumnHeader } from "@/components/ui/data-table/data-table-column-header";
import { toTitleCase } from "@/lib/utils";
import { getRegion, Group, Player } from "@/types";
import { TEAM_SIZE } from "@/types/constants";
import { Link } from "@tanstack/react-router";
import { Eye, EyeOff } from "lucide-react";

const defaultFilterFn = (
  row: Row<Group>,
  accessorKey: string,
  value: string,
) => {
  return value.includes(row.getValue(accessorKey));
};

export const columns: ColumnDef<Group>[] = [
  {
    accessorKey: "name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Name" />
    ),
    cell: ({ row }) => {
      const open = row.original.open;
      return (
        <div className="flex space-x-2">
          {open ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
          <span className="max-w-[500px] truncate font-medium">
            <Link to={`/groups/${row.original.id}`} className="hover:underline">
              {row.getValue("name")}
            </Link>
          </span>
        </div>
      );
    },
    enableHiding: false,
  },
  {
    accessorKey: "region",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Region" />
    ),
    cell: ({ row }) => {
      const region: string = row.getValue("region");

      if (!region) {
        return null;
      }

      return (
        <div className="flex w-[100px] items-center">
          <span>{getRegion(region)}</span>
        </div>
      );
    },
    filterFn: (row, accessorKey, value) => {
      const cellValue: string = row.getValue(accessorKey);
      if (!cellValue) {
        return false;
      }
      return value.includes(cellValue.toUpperCase());
    },
    enableHiding: false,
  },
  {
    accessorKey: "gamemode",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Gamemode" />
    ),
    cell: ({ row }) => {
      const gamemode: string = row.getValue("gamemode");

      if (!gamemode) {
        return null;
      }

      return (
        <div className="flex w-[100px] items-center">
          <span>{toTitleCase(gamemode)}</span>
        </div>
      );
    },
    filterFn: defaultFilterFn,
    enableHiding: false,
  },
  {
    accessorKey: "requirements",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Requirements" />
    ),
    cell: () => {
      return null;
    },
    enableHiding: false,
  },
  {
    accessorKey: "players",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Size" />
    ),
    cell: ({ row }) => {
      const players: Player[] = row.getValue("players");
      if (!players) {
        return null;
      }

      const playerCount = players.length;
      const innerText =
        playerCount === TEAM_SIZE ? "Full" : `${players.length}/${TEAM_SIZE}`;
      return (
        <div className="flex items-center">
          <span className="text-muted-foreground">{innerText}</span>
        </div>
      );
    },
    enableHiding: false,
  },
];
