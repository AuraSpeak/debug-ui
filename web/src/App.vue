<script setup lang="ts">
import { computed, ref } from "vue";
import { useStringWs } from "@/composables/useStringWs";
import WsLog from "@/components/WsLog.vue";
import WsSendInput from "@/components/WsSendInput.vue";
import ServerState from "./components/states/ServerState.vue";
import { useServerStore } from "@/stores/UDPServer";
import { useServerButton } from "./composables/useServerButton";
import ClientsState from "./components/states/ClientsState.vue";
import Overlay from "@/components/Overlay.vue";
import { useApi } from "@/api/useApi";

const udpServerButtonConfig = useServerButton();
const serverStore = useServerStore();
serverStore.fetchState();

const api = useApi();

const buttonConfig = udpServerButtonConfig.buttonConfig;

async function handleServerAction() {
  const action = buttonConfig.value.action;
  if (action === "connect") {
    await serverStore.startServer();
  } else if (action === "disconnect") {
    await serverStore.stopServer();
  }
}

type UsUEvent = { id: number, seq: number };
const NewClient = ref(false);
const usuEvent = ref<UsUEvent | null>(null);
let seq = 0;
const wsUrl = import.meta.env.DEV
  ? "ws://localhost:8080/ws"
  : (location.protocol === "https:" ? "wss://" : "ws://") + location.host + "/ws";

const { status, lines, error, send, connect, close } = useStringWs(wsUrl, {
  onMessage: (data) => {
    // Spezielle Behandlung für "uss" Nachricht (Update Server State)
    if (data === "uss") {
      serverStore.fetchState();
    } else if (data === "cnu") {
      console.log("cnu");
      NewClient.value = true;
    } else if (data.startsWith("usu,")) {
      console.log("usu", data);
      const parts = data.split(",");
      if (parts.length > 1 && parts[1]) {
        const clientId = parseInt(parts[1], 10);
        usuEvent.value = { id: clientId, seq: seq++ };
      }
    } else if (data === "rp") {
      console.log("rp");
      send("ack/rp");
      location.reload();
    }
  },
});

const canSend = computed(() => status.value === "open");

function handleSend(text: string) {
  const ok = send(text);
  if (!ok) lines.value.push("[ws] not connected");
}

// clients overlay
const clientsOverlay = ref(false);
</script>

<template>
  <div class="drawer">
    <input id="drawer-sidebar" type="checkbox" class="drawer-toggle">
    <div class="drawer-content">
      <label for="drawer-sidebar" class="btn drawer-button">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"
          stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </label>
      <div class="p-6 space-y-4">
        <div class="flex items-center gap-3">
          <span class="badge"
            :class="status === 'open' ? 'badge-success' : status === 'connecting' ? 'badge-warning' : 'badge-ghost'">
            {{ status }}
          </span>
          <span v-if="error" class="text-error text-sm">{{ error }}</span>

          <div class="ml-auto flex gap-2">
            <button class="btn btn-primary" @click="api.udpClients.start()">
              Start Client
            </button>
            <button :class="buttonConfig.class" :disabled="buttonConfig.disabled" @click="handleServerAction">
              {{ buttonConfig.label }}
            </button>
            <button class="btn btn-sm" @click="connect">Connect</button>
            <button class="btn btn-sm btn-outline" @click="close">Close</button>
          </div>
        </div>

        <WsLog :lines="lines" />

        <WsSendInput :disabled="!canSend" @send="handleSend" />
      </div>
    </div>
    <div class="drawer-side">
      <label for="drawer-sidebar" class="drawer-overlay"></label>
      <ul class="menu bg-base-200 min-h-full w-80 p-4">
        <li><button onclick="server_modal.showModal()">Server</button></li>
        <li><button @click="clientsOverlay = true">Clients</button></li>
      </ul>
    </div>
    <!-- Server Modal -->
    <dialog id="server_modal" class="modal">
      <div class="modal-box">
        <ServerState />
        <div class="modal-action">
          <form method="dialog">
            <button :class="buttonConfig.class" :disabled="buttonConfig.disabled" @click="handleServerAction">
              {{ buttonConfig.label }}
            </button>
            <button class="btn">Close</button>
          </form>
        </div>
      </div>
    </dialog>
    <!-- Clients Overlay -->
    <Overlay v-model="clientsOverlay" title="Clients" widthClass="w-11/12 w-[90vw]" maxWClass="max-w-none"
      heightClass="h-11/12">
      <ClientsState :needs-update="NewClient" :usu-event="usuEvent" @done:clients-update="NewClient = false" />
    </Overlay>
  </div>
</template>
