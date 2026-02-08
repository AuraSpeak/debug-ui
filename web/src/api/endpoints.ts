import type { ApiClient } from "./client";
import type { ID, Paginated, ServerState, UDPClient, UDPClientState, SendDatagramRequest, MermaidTraces, ClientMapData } from "./types";

export interface ServerApi {
    start: () => Promise<void>;
    stop: () => Promise<void>;
    getState: () => Promise<ServerState>;
}

export interface UDPClientApi {
    start: () => Promise<void>;
    stop: () => Promise<void>;
    getStateByName: (name: string) => Promise<UDPClientState>;
    getStateById: (id: ID) => Promise<UDPClientState>;
    getAll: () => Promise<{ udpClients: UDPClient[] }>;
    list: (params: { page?: number, pageSize?: number, q?: string }) => Promise<Paginated<UDPClient>>;
    sendDatagram: (request: SendDatagramRequest) => Promise<void>;
    getClientMap: () => Promise<ClientMapData>;
}

export interface TraceApi {
    getAll: (name: string) => Promise<MermaidTraces>;
}

export function createServerApi(client: ApiClient): ServerApi {
    return {
        start: () => client.post("/api/server/start"),
        stop: () => client.post("/api/server/stop"),
        getState: () => client.get("/api/server/get"),
    }
}

export function createUDPClientApi(client: ApiClient): UDPClientApi {
    return {
        start: () => client.post("/api/client/start"),
        stop: () => client.post("/api/client/stop"),
        getStateByName: (name: string) => client.get("/api/client/get/name", { query: { name } }),
        getStateById: (id: ID) => client.get("/api/client/get/id", { query: { id } }),
        getAll: () => client.get("/api/client/get/all"),
        list: (params: { page?: number, pageSize?: number, q?: string }) => client.get("/api/client/get/all/paginated", { query: params }),
        sendDatagram: (request: SendDatagramRequest) => client.post("/api/client/send", { body: request }),
        getClientMap: () => client.get<ClientMapData>("/api/client/map"),
    };
}

export function createTraceApi(client: ApiClient): TraceApi {
    return {
        getAll: (name: string) => client.get("/api/traces/all", { query: {name}}),
    };
}
