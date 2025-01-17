import { useEffect, useRef, useState } from "react";

import { useParams } from "@tanstack/react-router";
import { Send } from "lucide-react";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
  Skeleton,
} from "@/components/ui";
import { Button } from "@/components/ui/button";
import {
  ChatBubble,
  ChatBubbleMessage,
  ChatInput,
  ChatMessageList,
} from "@/components/ui/chat";
import { useGroupChat, useProfile } from "@/hooks";
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

  const [newMessage, setNewMessage] = useState("");
  const messagesRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (newMessage.trim()) {
      sendMessage(newMessage.trim());
      setNewMessage("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault(); // Prevents adding a new line in the input
      if (newMessage.trim()) {
        handleSendMessage(e as unknown as React.FormEvent);
      }
    }
  };

  const isConnected = connectionStatus === "connected";

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
              <ChatBubble key={message.id} variant="received">
                <ChatBubbleMessage>
                  <div className="flex items-center gap-2">
                    <span
                      className={`font-semibold ${getUserColor(message.sender)}`}
                    >
                      {message.sender}
                    </span>
                    <span className="text-xs text-zinc-500 dark:text-zinc-400">
                      {formatTimestamp(message.timestamp)}
                    </span>
                  </div>
                  <p className="text-sm break-words">{message.content}</p>
                </ChatBubbleMessage>
              </ChatBubble>
            ))}
          </ChatMessageList>
        )}
      </CardContent>
      <CardFooter className="p-2 w-full">
        <form
          onSubmit={handleSendMessage}
          className="flex align-center gap-2 w-full m-2"
        >
          <ChatInput
            value={newMessage}
            onKeyDown={handleKeyDown}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Send message"
            className="flex-1 w-full"
            disabled={!isPlayerInGroup || !isConnected}
          />
          <Button
            type="submit"
            size="icon"
            disabled={!isPlayerInGroup || !isConnected}
          >
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </CardFooter>
    </Card>
  );
}
