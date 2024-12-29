import { ColumnDef, Row } from "@tanstack/react-table";

import { DataTableColumnHeader } from "@/components/ui/data-table/data-table-column-header";
import { toTitleCase } from "@/lib/utils";
import {
  areRequirementsMet,
  Gamemode,
  gamemodeEmojis,
  getRankFromRankVal,
  getRegion,
  getRequirements,
  Group,
  Player,
  Profile,
} from "@/types";
import { TEAM_SIZE } from "@/types/constants";
import { Link } from "@tanstack/react-router";
import { Eye, EyeOff } from "lucide-react";
import {
  HoverCardContent,
  HoverCard,
  HoverCardTrigger,
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "./ui";

const defaultFilterFn = (
  row: Row<Group>,
  accessorKey: string,
  value: string,
) => {
  return value.includes(row.getValue(accessorKey));
};

export const columns = (profile: Profile): ColumnDef<Group>[] => [
  {
    accessorKey: "name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Name" />
    ),
    cell: ({ row }) => {
      const open = row.original.open;
      return (
        <div className="flex space-x-2">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>
                {open ? (
                  <Eye className="w-4 h-4" />
                ) : (
                  <EyeOff className="w-4 h-4" />
                )}
              </TooltipTrigger>
              <TooltipContent>{open ? "Public" : "Private"}</TooltipContent>
            </Tooltip>
          </TooltipProvider>
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
      const gamemode: Gamemode = row.getValue("gamemode");

      if (!gamemode) {
        return null;
      }

      return (
        <div className="flex w-[100px] items-center">
          <span>{`${gamemodeEmojis[gamemode]} ${toTitleCase(gamemode)}`}</span>
        </div>
      );
    },
    filterFn: defaultFilterFn,
    enableHiding: false,
  },
  {
    accessorKey: "requirements",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Requirements Met" />
    ),
    cell: ({ row }) => {
      const group = row.original as Group;
      const requirements = getRequirements(group);
      return (
        <div className="flex items-center">
          <HoverCard>
            <HoverCardTrigger asChild>
              <span>
                {areRequirementsMet(group, requirements, profile) ? "✅" : "❌"}
              </span>
            </HoverCardTrigger>
            <HoverCardContent className="w-full p-2">
              <ul>
                {requirements.voiceChat && (
                  <li>
                    <strong>Voice Chat:</strong> ✅
                  </li>
                )}
                {requirements.mic && (
                  <li>
                    <strong>Mic:</strong> ✅
                  </li>
                )}
                {group.gamemode !== Gamemode.Quickplay && group.roleQueue && (
                  <>
                    <li>
                      <strong>Min Rank:</strong>{" "}
                      {getRankFromRankVal(requirements.minRank)}
                    </li>
                    <li>
                      <strong>Max Rank:</strong>{" "}
                      {getRankFromRankVal(requirements.maxRank)}
                    </li>
                    <li>
                      <strong>Vanguards:</strong>{" "}
                      {requirements.requestedRoles.vanguards.curr}/
                      {requirements.requestedRoles.vanguards.max}
                    </li>
                    <li>
                      <strong>Duelists:</strong>{" "}
                      {requirements.requestedRoles.duelists.curr}/
                      {requirements.requestedRoles.duelists.max}
                    </li>
                    <li>
                      <strong>Strategists:</strong>{" "}
                      {requirements.requestedRoles.strategists.curr}/
                      {requirements.requestedRoles.strategists.max}
                    </li>
                  </>
                )}
              </ul>
            </HoverCardContent>
          </HoverCard>
        </div>
      );
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
      return (
        <div className="flex items-center">
          <span className="text-muted-foreground">{`${players.length}/${TEAM_SIZE}`}</span>
        </div>
      );
    },
    enableHiding: false,
  },
];
