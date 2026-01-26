import { inject } from "vue";
import { ApiKey } from "./index";

export function useApi() {
    const api = inject(ApiKey);
    if (!api) throw new Error("API not provided. Did you call installApi() in main.ts?");
    return api;
}