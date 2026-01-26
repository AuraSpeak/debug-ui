import { defineStore } from "pinia"
import { createApi } from "@/api"
import type { ServerState } from "@/api/types"

// Create API instance for use in the store
const api = createApi(import.meta.env.VITE_API_URL || "http://localhost:8080")

export const useServerStore = defineStore("server", {
  state: () => ({
    ServerState: null as ServerState | null,
    lastUpdated: null as number | null,

    // optional UX bits
    loading: false,
    error: null as string | null,
  }),

  getters: {
    status(state): "offline" | "running" | "stopping" | "unknown" {
      if (state.lastUpdated === null) return "unknown"
      if (!state.ServerState?.isAlive) return "offline"
      if (state.ServerState?.shouldStop) return "stopping"
      return "running"
    },
  },

  actions: {
    applyServerState(payload: ServerState) {
      this.ServerState = payload
      this.lastUpdated = Date.now()
    },

    async fetchState() {
      this.loading = true
      this.error = null
      try {
        const data = await api.server.getState()
        this.applyServerState(data)
      } catch (e: any) {
        this.error = e?.message ?? "Failed to fetch server state"
      } finally {
        this.loading = false
      }
    },

    async startServer() {
      this.loading = true
      this.error = null
      try {
        await api.server.start()
        await this.fetchState()
      } catch (e: any) {
        this.error = e?.message ?? "Failed to start server"
      } finally {
        this.loading = false
      }
    },

    async stopServer() {
      this.loading = true
      this.error = null
      try {
        await api.server.stop()
        await this.fetchState()
      } catch (e: any) {
        this.error = e?.message ?? "Failed to stop server"
      } finally {
        this.loading = false
      }
    },
  },
})