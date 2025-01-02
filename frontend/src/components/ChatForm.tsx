import { useEffect, useRef, useState } from "react";
import { Send } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";

interface Message {
  id: string;
  sender: string;
  content: string;
  timestamp: string;
}

function getTimeString() {
  const now = new Date();
  return now.toLocaleTimeString("en-US", {
    hour12: false,
    hour: "2-digit",
    minute: "2-digit",
  });
}

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

interface ChatFormProps {
  username: string;
}

export function ChatForm({ username = "User" }: ChatFormProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (newMessage.trim()) {
      const message: Message = {
        id: crypto.randomUUID(),
        sender: username,
        content: newMessage.trim(),
        timestamp: getTimeString(),
      };
      setMessages([...messages, message]);
      setNewMessage("");
    }
  };

  return (
    <div className="flex flex-col h-full">
      <ScrollArea className="flex-1 pr-4">
        <div className="space-y-4 pb-4">
          {messages.map((message) => (
            <div key={message.id} className="flex flex-col">
              <div className="flex items-center gap-2">
                <span
                  className={`font-semibold ${getUserColor(message.sender)}`}
                >
                  {message.sender}
                </span>
                <span className="text-xs text-zinc-500 dark:text-zinc-400">
                  {message.timestamp}
                </span>
              </div>
              <p className="text-sm">{message.content}</p>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
      </ScrollArea>

      <form onSubmit={handleSendMessage} className="border-t pt-4 flex gap-2">
        <Input
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type a message..."
          className="flex-1"
        />
        <Button type="submit" size="icon">
          <Send className="h-4 w-4" />
        </Button>
      </form>
    </div>
  );
}
