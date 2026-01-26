import { onBeforeUnmount, ref, type Ref } from "vue";

export type WebSocketStatus = "connecting" | "open" | "closed" | "error";

export interface UseWebSocketOptions {
  /** Automatisch beim Erstellen verbinden (Standard: true) */
  autoConnect?: boolean;
  /** Nachrichten in einem Array speichern (Standard: true) */
  storeMessages?: boolean;
  /** Callback wenn die Verbindung geöffnet wird */
  onOpen?: () => void;
  /** Callback wenn eine Nachricht empfangen wird */
  onMessage?: (data: string) => void;
  /** Callback wenn ein Fehler auftritt */
  onError?: (error: Event) => void;
  /** Callback wenn die Verbindung geschlossen wird */
  onClose?: (event: CloseEvent) => void;
  /** Nachrichten transformieren bevor sie gespeichert/übergeben werden */
  transformMessage?: (data: string) => string;
}

export interface UseWebSocketReturn {
  /** Aktueller Verbindungsstatus */
  status: Ref<WebSocketStatus>;
  /** Array aller empfangenen Nachrichten (nur wenn storeMessages: true) */
  lines: Ref<string[]>;
  /** Fehlermeldung falls vorhanden */
  error: Ref<string | null>;
  /** Nachricht senden */
  send: (text: string) => boolean;
  /** Verbindung herstellen */
  connect: () => void;
  /** Verbindung schließen */
  close: () => void;
  /** WebSocket-Instanz (für erweiterte Nutzung) */
  socket: Ref<WebSocket | null>;
}

/**
 * Vue Composable für WebSocket-Verbindungen mit klarer, konfigurierbarer API
 * 
 * @example
 * ```ts
 * const { status, lines, send, connect } = useStringWs("ws://localhost:8080/ws", {
 *   onMessage: (data) => console.log("Received:", data),
 *   onOpen: () => console.log("Connected"),
 * });
 * ```
 */
export function useStringWs(
  url: string | Ref<string>,
  options: UseWebSocketOptions = {}
): UseWebSocketReturn {
  const {
    autoConnect = true,
    storeMessages = true,
    onOpen,
    onMessage,
    onError,
    onClose,
    transformMessage,
  } = options;

  const status = ref<WebSocketStatus>("closed");
  const lines = ref<string[]>([]);
  const error = ref<string | null>(null);
  const socket = ref<WebSocket | null>(null);

  function connect() {
    // Schließe bestehende Verbindung falls vorhanden
    if (socket.value) {
      socket.value.close();
    }

    status.value = "connecting";
    error.value = null;

    const urlValue = typeof url === "string" ? url : url.value;
    const ws = new WebSocket(urlValue);
    socket.value = ws;

    ws.onopen = () => {
      status.value = "open";
      if (storeMessages) {
        lines.value.push("[ws] connected");
      }
      onOpen?.();
    };

    ws.onmessage = (ev) => {
      const rawData = typeof ev.data === "string" ? ev.data : String(ev.data);
      const processedData = transformMessage ? transformMessage(rawData) : rawData;

      if (storeMessages) {
        lines.value.push(processedData);
      }
      onMessage?.(processedData);
    };

    ws.onerror = (ev) => {
      status.value = "error";
      error.value = "WebSocket error";
      if (storeMessages) {
        lines.value.push("[ws] error");
      }
      onError?.(ev);
    };

    ws.onclose = (ev) => {
      setTimeout(() => {
        connect();
      }, 1000);
      status.value = "closed";
      if (storeMessages) {
        lines.value.push("[ws] closed");
      }
      onClose?.(ev);
      socket.value = null;
    };
  }

  function send(text: string): boolean {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
      return false;
    }
    socket.value.send(text);
    return true;
  }

  function close() {
    if (socket.value) {
      socket.value.close();
      socket.value = null;
    }
  }

  if (autoConnect) {
    connect();
  }

  onBeforeUnmount(close);

  return {
    status,
    lines,
    error,
    send,
    connect,
    close,
    socket,
  };
}
