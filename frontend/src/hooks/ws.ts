import { useCallback, useEffect, useRef, useState } from "react";

import { WebSocketClient, WebSocketMessage } from "@/api/ws";
import { useProfile } from "@/hooks";

export type WebSocketHandler = (message: WebSocketMessage) => void;

export function useWebSocket(groupId: string) {
  const [client, setClient] = useState<WebSocketClient | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<
    "connected" | "disconnected"
  >("disconnected");
  const [profile] = useProfile();

  const clientRef = useRef<WebSocketClient | null>(null);

  useEffect(() => {
    if (!groupId || !profile.name) return;

    const wsClient = new WebSocketClient(groupId);

    clientRef.current = wsClient;
    setClient(wsClient);
    wsClient.connect();
    setConnectionStatus("connected");

    return () => {
      wsClient.disconnect();
      clientRef.current = null;
      setClient(null);
      setConnectionStatus("disconnected");
    };
  }, [groupId, profile.name]);

  const subscribe = useCallback((handler: WebSocketHandler) => {
    if (!clientRef.current) return () => {};
    return clientRef.current.onMessage(handler);
  }, []);

  const send = useCallback(
    (op: number, payload: unknown) => {
      if (!clientRef.current || !profile.name) return;
      clientRef.current.sendMessage(op, payload);
    },
    [profile.name],
  );

  return {
    send,
    subscribe,
    connectionStatus,
    isConnected: client?.isConnected() ?? false,
  };
}
