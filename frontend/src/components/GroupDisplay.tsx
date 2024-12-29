import { getPlatform, getRank, getRegion, Group, Player } from "@/types";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableCaption,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";
import { strArrayToTitleCase, toTitleCase } from "@/lib/utils";
import { useMemo } from "react";

import teamUps from "@/assets/teamups.json";
import { TeamUpItem } from "./TeamUp";
import { Button } from "./ui/button";
import { Copy, X } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

interface GroupDisplayProps {
  group: Group;
}

export function GroupDisplay({ group }: GroupDisplayProps) {
  const { toast } = useToast();

  const leader = group.players.find((player) => player.leader);

  const { currVanguards, currDuelists, currStrategists, currCharacters } =
    useMemo(() => {
      return group.players.reduce(
        (acc, player) => {
          acc.currVanguards += player.roles.includes("vanguard") ? 1 : 0;
          acc.currDuelists += player.roles.includes("duelist") ? 1 : 0;
          acc.currStrategists += player.roles.includes("strategist") ? 1 : 0;
          acc.currCharacters = acc.currCharacters.union(
            new Set(player.characters),
          );
          return acc;
        },
        {
          currVanguards: 0,
          currDuelists: 0,
          currStrategists: 0,
          currCharacters: new Set(),
        },
      );
    }, [group.players]);

  // TODO: Exclude suggested teamups that would violate the role queue restrictions
  const suggestedTeamUps = useMemo(() => {
    return teamUps.filter(
      (teamup) =>
        new Set(teamup.requirements.allOf)
          .union(new Set(teamup.requirements.oneOf))
          .intersection(currCharacters).size > 0,
    );
  }, [teamUps, currCharacters]);

  function onKick(player: Player) {
    console.log(player);
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>
          {leader?.name}&apos;s Group
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
              <TableHead>Kick</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {group.players.map((player) => (
              <TableRow key={player.name}>
                <TableCell>
                  {player.name}
                  {player.leader ? " 🚩" : ""}
                </TableCell>
                <TableCell>{getRank(player.rank)}</TableCell>
                <TableCell>{strArrayToTitleCase(player.roles)}</TableCell>
                <TableCell>{player.characters.join(", ")}</TableCell>
                <TableCell>{getPlatform(player.platform)}</TableCell>
                <TableCell>
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => onKick(player)}
                  >
                    <X className="w-4 h-4" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
