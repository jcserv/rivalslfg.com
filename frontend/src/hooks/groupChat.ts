import { useCallback, useEffect, useState } from "react";

import { ChatMessage, WebSocketMessage, WebSocketOp } from "@/api/ws";
import { useProfile } from "@/hooks";

import { useWebSocket } from "./ws";

export function useGroupChat(groupId: string) {
  const ws = useWebSocket(groupId);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [profile] = useProfile();

  const messageHandler = useCallback((message: WebSocketMessage) => {
    switch (message.op) {
      case WebSocketOp.GroupChat: {
        const chatMessage = message.payload as ChatMessage;
        setMessages((prev) => [...prev, chatMessage]);
        break;
      }
      case WebSocketOp.GroupJoin:
        // Handle member join
        break;
      case WebSocketOp.GroupLeave:
        // Handle member leave
        break;
      case WebSocketOp.GroupPromotion:
        // Handle member promotion
        break;
    }
  }, []);

  useEffect(() => {
    const cleanup = ws.subscribe(messageHandler);
    return () => {
      cleanup();
      setMessages([]);
    };
  }, [ws.subscribe, messageHandler]);

  const sendMessage = useCallback(
    (content: string) => {
      if (!profile.name) return;
      ws.send(WebSocketOp.GroupChat, {
        id: crypto.randomUUID(),
        content,
        sender: profile.name,
        timestamp: new Date().toISOString(),
      });
    },
    [ws.send, profile.name],
  );

  return {
    messages,
    sendMessage,
    connectionStatus: ws.connectionStatus,
  };
}
