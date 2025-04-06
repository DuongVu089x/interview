export type Notification = {
  topic: string;
  title: string;
  description: string;
  link: string;
}

export type WebSocketMessage = {
  topic: string;
  content: Notification | Record<string, unknown>;
}

type MessageCallback = (data: WebSocketMessage) => void;

class WebSocketClient {
  private socket: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private messageListeners: MessageCallback[] = [];
  private connectListeners: (() => void)[] = [];
  private disconnectListeners: (() => void)[] = [];
  private errorListeners: ((error: Event) => void)[] = [];
  private url: string;

  constructor(url: string = 'ws://localhost:8382/notifications') {
    this.url = url;
  }

  connect() {
    if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
      return;
    }

    this.socket = new WebSocket(this.url);

    this.socket.onopen = () => {
      console.log('Connected to WebSocket server');
      this.reconnectAttempts = 0;
      this.connectListeners.forEach(listener => listener());
    };

    this.socket.onclose = () => {
      console.log('Disconnected from WebSocket server');
      this.disconnectListeners.forEach(listener => listener());

      // Attempt to reconnect
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        setTimeout(() => this.connect(), this.reconnectDelay * this.reconnectAttempts);
      }
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.errorListeners.forEach(listener => listener(error));
    };

    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.messageListeners.forEach(listener => listener(data));
      } catch (e) {
        console.error('Error parsing WebSocket message:', e);
      }
    };
  }

  send(message: WebSocketMessage) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error('Cannot send message: WebSocket is not connected');
    }
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  onMessage(callback: MessageCallback) {
    this.messageListeners.push(callback);
    return () => {
      this.messageListeners = this.messageListeners.filter(listener => listener !== callback);
    };
  }

  onConnect(callback: () => void) {
    this.connectListeners.push(callback);
    return () => {
      this.connectListeners = this.connectListeners.filter(listener => listener !== callback);
    };
  }

  onDisconnect(callback: () => void) {
    this.disconnectListeners.push(callback);
    return () => {
      this.disconnectListeners = this.disconnectListeners.filter(listener => listener !== callback);
    };
  }

  onError(callback: (error: Event) => void) {
    this.errorListeners.push(callback);
    return () => {
      this.errorListeners = this.errorListeners.filter(listener => listener !== callback);
    };
  }

  isConnected() {
    return this.socket && this.socket.readyState === WebSocket.OPEN;
  }
}

// Singleton instance
let wsClient: WebSocketClient | null = null;

export const initializeSocket = (url?: string) => {
  if (!wsClient) {
    wsClient = new WebSocketClient(url);
    wsClient.connect();
  }
  return wsClient;
};

export const getSocket = () => {
  if (!wsClient) {
    throw new Error('WebSocket not initialized. Call initializeSocket first.');
  }
  return wsClient;
};

export const closeSocket = () => {
  if (wsClient) {
    wsClient.disconnect();
    wsClient = null;
  }
};
