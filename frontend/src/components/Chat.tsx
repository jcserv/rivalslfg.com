import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Skeleton,
} from "@/components/ui";

interface ChatProps {
  canUserAccessGroup: boolean | null;
}

export function Chat({ canUserAccessGroup }: ChatProps) {
  return (
    <Card className="h-1/2">
      <CardHeader>
        <CardTitle>Chat</CardTitle>
      </CardHeader>
      <CardContent className="h-full">
        {!canUserAccessGroup && (
          <Skeleton className="h-3/4 w-full rounded-xl" />
        )}
        {canUserAccessGroup && <p>Ping!</p>}
      </CardContent>
    </Card>
  );
}
