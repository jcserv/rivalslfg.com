export const WebSocketOp = {
  GroupChat: 1,
  GroupJoin: 2,
  GroupLeave: 3,
  GroupPromotion: 4,
} as const;

export type WebSocketMessage = {
  groupId: string;
  op: number;
  payload: unknown;
};

export type ChatMessage = {
  id: string;
  sender: string;
  content: string;
  timestamp: string;
};

type MessageHandler = (message: WebSocketMessage) => void;

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private messageHandlers: Set<MessageHandler> = new Set();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1s delay

  constructor(
    private groupId: string,
    private baseUrl?: string,
    onMessage?: MessageHandler,
  ) {
    if (onMessage) {
      this.messageHandlers.add(onMessage);
    }

    // Convert http(s) to ws(s)
    this.baseUrl = import.meta.env.VITE_API_URL.replace(/^http/, "ws");
  }

  private getToken(): string | null {
    return localStorage.getItem("token");
  }

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      const wsUrl = `${this.baseUrl}/ws?groupId=${this.groupId}&access_token=${this.getToken()}`;

      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000;
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage;
          this.messageHandlers.forEach((handler) => handler(message));
        } catch {
          // console.error("Error parsing WebSocket message:", error);
        }
      };

      this.ws.onclose = () => {
        this.attemptReconnect();
      };

      this.ws.onerror = () => {
        // console.error("WebSocket error:", error);
      };
    } catch {
      // console.error("Error connecting to WebSocket:", error);
      this.attemptReconnect();
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(
      this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1),
      30000,
    ); // Max 30s

    setTimeout(() => {
      this.connect();
    }, delay);
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  sendMessage(op: number, payload: unknown) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      return;
    }

    const message: WebSocketMessage = {
      groupId: this.groupId,
      op,
      payload,
    };

    try {
      this.ws.send(JSON.stringify(message));
    } catch {
      //console.error("Error sending WebSocket message:", error);
    }
  }

  sendChatMessage(content: string, sender: string) {
    this.sendMessage(WebSocketOp.GroupChat, {
      id: crypto.randomUUID(),
      content,
      sender,
      timestamp: new Date().toISOString(),
    });
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => {
      this.messageHandlers.delete(handler);
    };
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}
