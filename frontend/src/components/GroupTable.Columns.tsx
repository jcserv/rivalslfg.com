import { ColumnDef, Row } from "@tanstack/react-table";

import { DataTableColumnHeader } from "@/components/ui/data-table/data-table-column-header";
import { toTitleCase } from "@/lib/utils";
import {
  Gamemode,
  gamemodeEmojis,
  getRankFromRankVal,
  getRegion,
  Group,
  GroupRequirements,
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

export type GroupTableData = Group & {
  requirements: GroupRequirements;
  areRequirementsMet: boolean;
};

const defaultFilterFn = (
  row: Row<GroupTableData>,
  accessorKey: string,
  value: string,
) => {
  return value.includes(row.getValue(accessorKey));
};

export const columns = (
  profile: Profile | null,
  isProfileEmpty: boolean,
): ColumnDef<GroupTableData>[] => [
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
            <Link
              to={`/groups/${row.original.id}`}
              className={!isProfileEmpty ? "hover:underline" : ""}
              disabled={isProfileEmpty}
            >
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
    accessorKey: "areRequirementsMet",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title="Requirements Met"
        className="justify-center"
      />
    ),
    cell: ({ row }) => {
      const group = row.original as GroupTableData;
      const areRequirementsMet = row.getValue("areRequirementsMet");
      const requirements: GroupRequirements = row.original.requirements;
      return (
        <div className="flex items-center text-center justify-center">
          <HoverCard>
            <HoverCardTrigger asChild>
              <span>{areRequirementsMet ? "✅" : "❌"}</span>
            </HoverCardTrigger>
            <HoverCardContent className="w-full p-2">
              <ul>
                {requirements.voiceChat && (
                  <li>
                    <strong>Voice Chat:</strong> Required
                  </li>
                )}
                {requirements.mic && (
                  <li>
                    <strong>Mic:</strong> Required
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
                  </>
                )}
              </ul>
            </HoverCardContent>
          </HoverCard>
        </div>
      );
    },
    filterFn: (
      row: Row<GroupTableData>,
      accessorKey: string,
      value: string[],
    ) => {
      console.log(value, accessorKey, row.getValue(accessorKey));
      const compareWith = row.getValue(accessorKey);
      return value.some((str) => {
        const boolValue =
          str.toLowerCase() === "true" ||
          str === "1" ||
          str.toLowerCase() === "yes";
        return boolValue === compareWith;
      });
    },
    enableHiding: false,
  },
  {
    accessorKey: "players",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title="Team"
        className="justify-center"
      />
    ),
    cell: ({ row }) => {
      const group = row.original as GroupTableData;
      const players: Player[] = row.getValue("players");
      if (!players) {
        return null;
      }
      const requirements: GroupRequirements = row.original.requirements;
      return (
        <div className="flex items-center text-center justify-center">
          <ul>
            {group.roleQueue && requirements.requestedRoles.vanguards.max && (
              <li>
                <strong>Vanguards:</strong>{" "}
                {requirements.requestedRoles.vanguards.curr}/
                {requirements.requestedRoles.vanguards.max}
              </li>
            )}
            {group.roleQueue && requirements.requestedRoles.duelists.max && (
              <li>
                <strong>Duelists:</strong>{" "}
                {requirements.requestedRoles.duelists.curr}/
                {requirements.requestedRoles.duelists.max}
              </li>
            )}
            {group.roleQueue && requirements.requestedRoles.strategists.max && (
              <li>
                <strong>Strategists:</strong>{" "}
                {requirements.requestedRoles.strategists.curr}/
                {requirements.requestedRoles.strategists.max}
              </li>
            )}
            <li>
              <span className="text-muted-foreground">
                Team: {`${players.length}/${TEAM_SIZE}`}
              </span>
            </li>
          </ul>
        </div>
      );
    },
    enableHiding: false,
  },
];
