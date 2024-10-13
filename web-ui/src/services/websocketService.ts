// websocketService.ts
type WebSocketAction = {
  action: string;
  payload: Record<string, any>;
};

export class WebSocketService {
  private static instance: WebSocketService | null = null;
  private urlBase = import.meta.env.VITE_WS_URL || "ws://localhost:8000";
  private url: string;
  private ws: WebSocket | null;
  private eventHandlers: Record<string, (message: WebSocketAction) => void>;

  // Singleton pattern to ensure a single instance of WebSocketService
  public static getInstance(url: string): WebSocketService {
    if (!WebSocketService.instance) {
      WebSocketService.instance = new WebSocketService(url);
    }
    return WebSocketService.instance;
  }

  // Private constructor to prevent direct instantiation
  private constructor(url: string) {
    this.url = this.urlBase + url;
    console.log("WebSocket URL:", this.url);
    this.ws = null;
    this.eventHandlers = {};
  }

  // Establish a WebSocket connection
  public connect(): void {
    if (this.ws) {
      console.warn("WebSocket is already connected.");
      return;
    }

    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log("WebSocket connected");
    };

    this.ws.onmessage = (event) => {
      const message: WebSocketAction = JSON.parse(event.data);
      console.log("New message received:", message);
      const { action } = message;
      if (this.eventHandlers[action]) {
        this.eventHandlers[action](message);
      }
    };

    this.ws.onclose = () => {
      console.log("WebSocket connection closed");
      this.ws = null;
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  // Send a message through the WebSocket
  public sendMessage(payload: WebSocketAction): void {
    if (!this.ws) {
      console.error("WebSocket is not connected. Cannot send message.");
      return;
    }

    const message = JSON.stringify(payload);
    this.ws.send(message);
  }

  // Register an event handler for a specific action
  public on(action: string, handler: (message: WebSocketAction) => void): void {
    this.eventHandlers[action] = handler;
  }

  // Close the WebSocket connection
  public close(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    } else {
      console.warn("WebSocket is already closed.");
    }
  }
}
