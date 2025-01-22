import { useEffect, useRef } from "react";

import { useParams } from "@tanstack/react-router";
import { Send } from "lucide-react";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
  Form,
  FormField,
  FormItem,
  FormMessage,
  Skeleton,
} from "@/components/ui";
import { Button } from "@/components/ui/button";
import {
  ChatBubble,
  ChatBubbleMessage,
  ChatInput,
  ChatMessageList,
} from "@/components/ui/chat";
import { ChatMessage, useChatForm, useGroupChat, useProfile } from "@/hooks";
import { formatTimestamp } from "@/lib";

const userColors = [
  "text-blue-500 dark:text-blue-400",
  "text-green-500 dark:text-green-400",
  "text-purple-500 dark:text-purple-400",
  "text-yellow-500 dark:text-yellow-400",
  "text-red-500 dark:text-red-400",
  "text-pink-500 dark:text-pink-400",
];

const userColorMap = new Map<string, string>();
let colorIndex = 0;

function getUserColor(username: string): string {
  if (!userColorMap.has(username)) {
    userColorMap.set(username, userColors[colorIndex % userColors.length]);
    colorIndex++;
  }
  return userColorMap.get(username) || userColors[0];
}

interface ChatBoxProps {
  canUserAccessGroup: boolean | null;
  isPlayerInGroup: boolean | undefined;
}

export function ChatBox({ canUserAccessGroup, isPlayerInGroup }: ChatBoxProps) {
  const { groupId } = useParams({ from: "/groups/$groupId" });
  const [profile] = useProfile();
  const { messages, sendMessage, connectionStatus } = useGroupChat(groupId);

  const messagesRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages]);

  const { form, handleSubmit, handleKeyDown } = useChatForm({
    onSubmit: (message) => sendMessage(message),
  });

  const isConnected = connectionStatus === "connected";
  const isDisabled = !isPlayerInGroup || !isConnected;
  const isSendDisabled = isDisabled || form.watch("message").length === 0;

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center">
          <span>Chat</span>
          <span
            className={`ml-2 text-sm ${isConnected ? "text-green-500" : "text-red-500"}`}
          >
            ●
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-0 h-[400px]">
        {" "}
        {!canUserAccessGroup && (
          <Skeleton className="h-full w-full rounded-xl" />
        )}
        {canUserAccessGroup && profile?.name && (
          <ChatMessageList className="flex-1" ref={messagesRef}>
            {messages.map((message) => (
              <ChatItem key={message.id} message={message} />
            ))}
          </ChatMessageList>
        )}
      </CardContent>
      <CardFooter className="p-2 w-full">
        <Form {...form}>
          <form
            onSubmit={handleSubmit}
            className="flex align-center gap-2 w-full m-2"
          >
            <FormField
              control={form.control}
              name="message"
              render={({ field }) => (
                <FormItem className="flex-1 w-full">
                  <ChatInput
                    id="message"
                    {...field}
                    onKeyDown={handleKeyDown}
                    placeholder="Send message"
                    disabled={isDisabled}
                  />
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" size="icon" disabled={isSendDisabled}>
              <Send className="h-4 w-4" />
            </Button>
          </form>
        </Form>
      </CardFooter>
    </Card>
  );
}

interface ChatMessageProps {
  message: ChatMessage;
}

function ChatItem({ message }: ChatMessageProps) {
  if (message.system) {
    return (
      <div>
        <span className="flex items-center justify-between">
          <p className="text-sm break-words italic mr-2">{message.content}</p>
          <p className="text-sm">{formatTimestamp(message.timestamp)}</p>
        </span>
      </div>
    );
  }
  return (
    <ChatBubble variant="received">
      <ChatBubbleMessage>
        <div className="flex items-center gap-2">
          <span className={`font-semibold ${getUserColor(message.sender)}`}>
            {message.sender}
          </span>
          <span className="text-xs text-zinc-500 dark:text-zinc-400">
            {formatTimestamp(message.timestamp)}
          </span>
        </div>
        <p className="text-sm break-words">{message.content}</p>
      </ChatBubbleMessage>
    </ChatBubble>
  );
}
