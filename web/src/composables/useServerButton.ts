import { computed } from "vue"
import { useServerStore } from "@/stores/UDPServer"

export function useServerButton() {
  const server = useServerStore()

  const buttonConfig = computed(() => {
    switch (server.status) {
      case "offline":
        return {
          label: "Connect",
          class: "btn btn-success",
          disabled: false,
          action: "connect" as const,
        }

      case "running":
        return {
          label: "Disconnect",
          class: "btn btn-error",
          disabled: false,
          action: "disconnect" as const,
        }

      case "stopping":
        return {
          label: "Stopping…",
          class: "btn btn-warning loading",
          disabled: true,
        }

      default:
        return {
          label: "Checking…",
          class: "btn btn-ghost loading",
          disabled: true,
        }
    }
  })

  return { buttonConfig }
}
