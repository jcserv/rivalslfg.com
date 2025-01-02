import { useMemo } from "react";

import { Copy, X } from "lucide-react";

import { TeamUpItem } from "@/components/TeamUp";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useProfile, useToast } from "@/hooks";
import { strArrayToTitleCase, toTitleCase } from "@/lib";
import { getGroupInfo, getPlatform, getRank, getRegion, Group } from "@/types";

import teamUps from "@/assets/teamups.json";

interface GroupDisplayProps {
  group: Group | undefined;
  canUserAccessGroup: boolean | null;
  isOwner: boolean;
  onRemove: (id: number, playerToRemove: string) => void;
}

export function GroupDisplay({
  group,
  canUserAccessGroup,
  isOwner,
  onRemove,
}: GroupDisplayProps) {
  const [profile] = useProfile();
  const { toast } = useToast();

  const { currVanguards, currDuelists, currStrategists, currCharacters } =
    useMemo(() => {
      return getGroupInfo(group);
    }, [group]);

  const suggestedTeamUps = useMemo(() => {
    return teamUps.filter(
      (teamup) =>
        new Set(teamup.requirements.allOf)
          .union(new Set(teamup.requirements.oneOf))
          .intersection(currCharacters).size > 0,
    );
  }, [teamUps, currCharacters]);

  if (!group) return null;
  return (
    <Card>
      <CardHeader>
        <CardTitle>
          {group?.name}
          <Button variant="outline" size="icon" className="ml-2">
            <Copy
              className="w-4 h-4"
              onClick={() => {
                navigator.clipboard.writeText(window.location.href);
                toast({
                  title: "Copied current URL to clipboard!",
                  variant: "success",
                });
              }}
            />
          </Button>
        </CardTitle>
        <CardDescription>
          <div className="grid grid-cols-12 gap-4">
            <div className="col-span-6">
              Region: {getRegion(group.region)}
              <br />
              Gamemode: {toTitleCase(group.gamemode)}
              <br />
              {group.roleQueue && (
                <>
                  Role Queue: {group.roleQueue ? "Enabled" : "Disabled"}
                  <br />
                  <br />
                  Team Composition:
                  <br />
                  <ul>
                    {group.roleQueue?.vanguards > 0 && (
                      <li>
                        •{" "}
                        {currVanguards < group.roleQueue.vanguards
                          ? "🌀"
                          : "✅"}{" "}
                        {currVanguards}/{group.roleQueue.vanguards} Vanguards{" "}
                      </li>
                    )}
                    {group.roleQueue?.duelists > 0 && (
                      <li>
                        •{" "}
                        {currDuelists < group.roleQueue.duelists ? "🌀" : "✅"}{" "}
                        {currDuelists}/{group.roleQueue.duelists} Duelists{" "}
                      </li>
                    )}
                    {group.roleQueue?.strategists > 0 && (
                      <li>
                        •{" "}
                        {currStrategists < group.roleQueue.strategists
                          ? "🌀"
                          : "✅"}{" "}
                        {currStrategists}/{group.roleQueue.strategists}{" "}
                        Strategists{" "}
                      </li>
                    )}
                  </ul>
                </>
              )}
            </div>
            {canUserAccessGroup && (
              <div className="col-span-6 flex justify-end">
                <div>
                  Suggested Team-ups:
                  <br />
                  <ul>
                    {suggestedTeamUps.map((teamup) => (
                      <TeamUpItem key={teamup.name} teamup={teamup} />
                    ))}
                  </ul>
                </div>
              </div>
            )}
          </div>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableCaption>Looking for players to team up with...</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>Player</TableHead>
              <TableHead>Rank</TableHead>
              <TableHead>Roles</TableHead>
              <TableHead>Characters</TableHead>
              <TableHead>Platform</TableHead>
              {isOwner && <TableHead>Kick</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {group.players.map((player) => (
              <TableRow key={player.name} isLoading={!canUserAccessGroup}>
                <TableCell>
                  {player.name}
                  {player.leader ? " 🚩" : ""}
                </TableCell>
                <TableCell>{getRank(player.rank)}</TableCell>
                <TableCell>{strArrayToTitleCase(player.roles)}</TableCell>
                <TableCell>{player.characters.join(", ")}</TableCell>
                <TableCell>{getPlatform(player.platform)}</TableCell>
                {isOwner && (
                  <TableCell>
                    {player.name !== profile.name && (
                      <Button
                        variant="outline"
                        size="icon"
                        onClick={() => onRemove(player.id, player.name)}
                      >
                        <X className="w-4 h-4" />
                      </Button>
                    )}
                  </TableCell>
                )}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
