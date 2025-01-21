import { useCallback, useEffect, useState } from "react";

import { useRouter } from "@tanstack/react-router";

import { WebSocketMessage, WebSocketOp } from "@/api/ws";
import {
  addPlayerToGroup,
  removePlayerFromGroup,
  useProfile,
  useToast,
} from "@/hooks";
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
  playerName: string;
  leaderId: number;
};

export function useGroupChat(groupId: string) {
  const router = useRouter();
  const { toast } = useToast();

  const ws = useWebSocket(groupId);
  const [profile] = useProfile();

  const [messages, setMessages] = useState<ChatMessage[]>([]);

  const messageHandler = useCallback(async (message: WebSocketMessage) => {
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
            content: `${payload.playerName} has left the group.`,
            timestamp: new Date().toISOString(),
          },
        ]);
        break;
      }
      case WebSocketOp.GroupDelete:
        toast({
          title: "Group was deleted",
          description: "You have been removed from the group.",
          variant: "destructive",
        });
        router.navigate({ to: ".." });
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
