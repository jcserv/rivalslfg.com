import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui";
import { formatSeasonBonus, TeamUp } from "@/types";

interface TeamUpProps {
  teamup: TeamUp;
}

export function TeamUpItem({ teamup }: TeamUpProps) {
  // TODO: Display the current & required characters for the teamup, red border if missing. green border if present.
  return (
    <li>
      <HoverCard>
        <HoverCardTrigger asChild>
          <span>
            •{" "}
            <p className="inline cursor-pointer hover:underline">
              {teamup.name}
            </p>
          </span>
        </HoverCardTrigger>
        <HoverCardContent className="w-80">
          <p>
            <strong>Seasonal Bonus -</strong>{" "}
            {formatSeasonBonus(teamup.seasonBonus)}
          </p>
          <p>
            <strong>Description -</strong> {teamup.description}
          </p>
        </HoverCardContent>
      </HoverCard>
    </li>
  );
}
