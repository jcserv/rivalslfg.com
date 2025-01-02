import { ChatForm } from "@/components";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Skeleton,
} from "@/components/ui";
import { useProfile } from "@/hooks";

interface ChatBoxProps {
  canUserAccessGroup: boolean | null;
}

export function ChatBox({ canUserAccessGroup }: ChatBoxProps) {
  const [profile] = useProfile();
  return (
    <Card className="h-full">
      <CardHeader className="pb-4">
        <CardTitle>Chat</CardTitle>
      </CardHeader>
      <CardContent className="p-4 h-[calc(100%-5rem)]">
        {!canUserAccessGroup && (
          <Skeleton className="h-full w-full rounded-xl" />
        )}
        {canUserAccessGroup && profile?.name && (
          <ChatForm username={profile.name} />
        )}
      </CardContent>
    </Card>
  );
}
