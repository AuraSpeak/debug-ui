// src/api/client.ts
import type { ApiErrorBody } from "./types";

export class ApiError extends Error {
  status: number;
  body?: ApiErrorBody;
  constructor(message: string, status: number, body?: ApiErrorBody) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }
}

export interface RequestOptions<TBody = unknown> {
  headers?: Record<string, string>;
  query?: Record<string, string | number | boolean | undefined | null>;
  body?: TBody;
  signal?: AbortSignal;
}

export interface ApiClient {
  get<T>(path: string, options?: RequestOptions): Promise<T>;
  post<T, TBody = unknown>(path: string, options?: RequestOptions<TBody>): Promise<T>;
  put<T, TBody = unknown>(path: string, options?: RequestOptions<TBody>): Promise<T>;
  patch<T, TBody = unknown>(path: string, options?: RequestOptions<TBody>): Promise<T>;
  del<T>(path: string, options?: RequestOptions): Promise<T>;
}

function buildUrl(baseUrl: string, path: string, query?: RequestOptions["query"]) {
  const url = new URL(path.replace(/^\//, ""), baseUrl.endsWith("/") ? baseUrl : baseUrl + "/");
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v === undefined || v === null) continue;
      url.searchParams.set(k, String(v));
    }
  }
  return url.toString();
}

async function parseJsonSafe(res: Response): Promise<any | undefined> {
  const text = await res.text();
  if (!text) return undefined;
  try {
    return JSON.parse(text);
  } catch {
    return undefined;
  }
}

export function createFetchClient(params: {
  baseUrl: string;
  getAuthToken?: () => string | null;
  defaultHeaders?: Record<string, string>;
}): ApiClient {
  const { baseUrl, getAuthToken, defaultHeaders } = params;

  async function request<T, TBody = unknown>(
    method: string,
    path: string,
    options: RequestOptions<TBody> = {},
  ): Promise<T> {
    const url = buildUrl(baseUrl, path, options.query);

    const token = getAuthToken?.();
    const headers: Record<string, string> = {
      ...(defaultHeaders ?? {}),
      ...(options.headers ?? {}),
    };

    // Only set JSON header when we actually send a body
    const hasBody = options.body !== undefined && options.body !== null;
    if (hasBody) headers["Content-Type"] = headers["Content-Type"] ?? "application/json";
    if (token) headers["Authorization"] = `Bearer ${token}`;

    const res = await fetch(url, {
      method,
      headers,
      body: hasBody ? JSON.stringify(options.body) : undefined,
      signal: options.signal,
    });

    if (!res.ok) {
      const body = (await parseJsonSafe(res)) as ApiErrorBody | undefined;
      const message = body?.message ?? `Request failed (${res.status})`;
      throw new ApiError(message, res.status, body);
    }

    // 204 No Content
    if (res.status === 204) return undefined as unknown as T;

    const data = await parseJsonSafe(res);
    return data as T;
  }

  return {
    get: (p, o) => request("GET", p, o),
    post: (p, o) => request("POST", p, o),
    put: (p, o) => request("PUT", p, o),
    patch: (p, o) => request("PATCH", p, o),
    del: (p, o) => request("DELETE", p, o),
  };
}
