export type ID = number;

export type ApiError = {
	code: number;
	message: string;
	details: string;
};

export type ApiSuccess = {
	message: string;
	details: string;
};

export interface ApiErrorBody {
    message: string;
    code?: number;
    details?: string;
}

export interface Paginated<T> {
    items: T[];
    page: number;
    pageSize: number;
    total: number;
}

export interface ServerState {
    shouldStop: boolean;
    isAlive: boolean;
}

export const DatagramDirection = {
    ClientToServer: 1,
    ServerToClient: 2,
} as const;

export type DatagramDirection = typeof DatagramDirection[keyof typeof DatagramDirection];

export interface Datagram {
    direction: typeof DatagramDirection[keyof typeof DatagramDirection];
    message: Uint8Array;
}

export interface UDPClient{
    id: ID;
    name: string;
}

export interface UDPClientState {
    id: ID;
    running: boolean;
    datagrams: Datagram[];
}

export interface SendDatagramRequest {
    id: ID;
    message: string;
    format: "hex" | "text";
}

export interface MermaidTraces {
    heading: string
    diagram: string
}

export interface LogEntry {
    level: string;
    msg: string;
    time: string;
    [key: string]: any; // Für zukünftige dynamische Felder
}

export interface ClientMapConnection {
    fromClientId: number;
    toClientId: number; // 0 = server
}

export interface ClientMapData {
    clients: Array<{ id: number; name: string }>;
    connections: ClientMapConnection[];
}
