import { getPlatform, getRank, getRegion, Group } from "@/types";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import {
  Table,
  TableCaption,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "./ui/table";
import { strArrayToTitleCase, toTitleCase } from "@/lib/utils";
import { useMemo } from "react";

interface GroupDisplayProps {
  group: Group;
}

export function GroupDisplay({ group }: GroupDisplayProps) {
  const leader = group.players.find((player) => player.leader);

  const { currVanguards, currDuelists, currStrategists } = useMemo(() => {
    return group.players.reduce(
      (acc, player) => {
        acc.currVanguards += player.roles.includes("vanguard") ? 1 : 0;
        acc.currDuelists += player.roles.includes("duelist") ? 1 : 0;
        acc.currStrategists += player.roles.includes("strategist") ? 1 : 0;
        return acc;
      },
      { currVanguards: 0, currDuelists: 0, currStrategists: 0 }
    );
  }, [group.players]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>{leader?.name}&apos;s Group</CardTitle>
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
                  <li>• Eye of Agabobo </li>
                  <li>• Unlimited Ammo </li>
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
              <TableHead>Ready</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {group.players.map((player) => (
                <TableRow key={player.name}>
                  <TableCell>{player.name}{player.leader ? " 🚩" : ""}</TableCell>
                  <TableCell>{getRank(player.rank)}</TableCell>
                  <TableCell>{strArrayToTitleCase(player.roles)}</TableCell>
                  <TableCell>{player.characters.join(", ")}</TableCell>
                  <TableCell>{getPlatform(player.platform)}</TableCell>
                  <TableCell>✅</TableCell>
                </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
