import { useCallback, useEffect, useState } from "react";

import { WebSocketMessage, WebSocketOp } from "@/api/ws";
import { addPlayerToGroup, removePlayerFromGroup, useProfile } from "@/hooks";
import { Player } from "@/types";

import { useWebSocket } from "./ws";

export type ChatMessage = {
  id: string;
  system?: boolean;
  sender: string;
  content: string;
  timestamp: string;
};

type PlayerLeftPayload = {
  playerId: number;
  leaderId: number;
};

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
      case WebSocketOp.GroupJoin: {
        const player = message.payload as Player;

        addPlayerToGroup(groupId, player);
        setMessages((prev) => [
          ...prev,
          {
            id: crypto.randomUUID(),
            system: true,
            sender: "System",
            content: `${player.name} has joined the group.`,
            timestamp: new Date().toISOString(),
          },
        ]);
        break;
      }
      case WebSocketOp.GroupLeave: {
        const payload = message.payload as PlayerLeftPayload;
        removePlayerFromGroup(groupId, payload.playerId, payload.leaderId);
        setMessages((prev) => [
          ...prev,
          {
            id: crypto.randomUUID(),
            system: true,
            sender: "System",
            content: `${payload.playerId} has left the group.`, // TODO: add player name
            timestamp: new Date().toISOString(),
          },
        ]);
        break;
      }
      case WebSocketOp.GroupDelete:
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
