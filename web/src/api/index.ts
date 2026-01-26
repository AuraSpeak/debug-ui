import type { App, InjectionKey} from "vue";
import { createFetchClient, type ApiClient } from "./client";
import { createServerApi, createTraceApi, createUDPClientApi, type ServerApi, type TraceApi, type UDPClientApi } from "./endpoints";

export interface Api {
    client: ApiClient;
    server: ServerApi;
    udpClients: UDPClientApi;
    trace: TraceApi;
}

export const ApiKey: InjectionKey<Api> = Symbol("Api");

export function createApi(baseUrl: string): Api {
    const client = createFetchClient({ 
        baseUrl,
        defaultHeaders: {
            "Content-Type": "application/json",
        },
    });
    return {
        client,
        server: createServerApi(client),
        udpClients: createUDPClientApi(client),
        trace: createTraceApi(client),
    }
}

export function installApi(app: App, api: Api) {
    app.provide(ApiKey, api);
}