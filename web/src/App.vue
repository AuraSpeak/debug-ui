<script setup lang="ts">
import { computed, ref, provide } from "vue";
import { useStringWs } from "@/composables/useStringWs";
import WsLog from "@/components/WsLog.vue";
import WsSendInput from "@/components/WsSendInput.vue";
import ServerState from "@/components/states/ServerState.vue";
import { useServerStore } from "@/stores/UDPServer";
import { useServerButton } from "@/composables/useServerButton";
import ClientsState from "@/components/states/ClientsState.vue";
import ClientMap from "@/components/ClientMap.vue";
import Overlay from "@/components/Overlay.vue";
import { useApi } from "@/api/useApi";

type Section = "log" | "server" | "clients" | "map";
const activeSection = ref<Section>("log");

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

type UsuEvent = { id: number; seq: number };
const newClient = ref(false);
const usuEvent = ref<UsuEvent | null>(null);
let seq = 0;

const mapRefreshTrigger = ref(0);
export type PacketEvent = { from: number; to: number; dir: number };
const packetEvent = ref<PacketEvent | null>(null);
provide("mapRefreshTrigger", mapRefreshTrigger);
provide("packetEvent", packetEvent);

const wsUrl = import.meta.env.DEV
  ? "ws://localhost:8080/ws"
  : (location.protocol === "https:" ? "wss://" : "ws://") + location.host + "/ws";

const { status, lines, error, send, connect, close } = useStringWs(wsUrl, {
  onMessage: (data) => {
    if (data === "uss") {
      serverStore.fetchState();
    } else if (data === "cnu") {
      newClient.value = true;
    } else if (data.startsWith("usu")) {
      const rest = data.slice(3);
      if (rest) {
        const clientId = parseInt(rest, 10);
        if (!Number.isNaN(clientId)) {
          usuEvent.value = { id: clientId, seq: seq++ };
        }
      }
    } else if (data === "map") {
      mapRefreshTrigger.value++;
    } else if (data.startsWith("pkt,")) {
      const parts = data.split(",");
      const p1 = parts[1];
      const p2 = parts[2];
      const p3 = parts[3];
      if (p1 !== undefined && p2 !== undefined && p3 !== undefined) {
        const from = parseInt(p1, 10);
        const to = parseInt(p2, 10);
        const dir = parseInt(p3, 10);
        if (!Number.isNaN(from) && !Number.isNaN(to) && !Number.isNaN(dir)) {
          packetEvent.value = { from, to, dir };
        }
      }
    } else if (data === "rp") {
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

const clientsOverlay = ref(false);
</script>

<template>
  <div class="min-h-screen flex flex-col bg-[var(--color-bg)]">
    <header class="border-b border-[var(--color-border)] bg-[var(--color-bg-elevated)] flex-shrink-0">
      <nav class="flex items-center gap-1 p-2">
        <button
          v-for="tab in [
            { id: 'log' as Section, label: 'Log' },
            { id: 'server' as Section, label: 'Server' },
            { id: 'clients' as Section, label: 'Clients' },
            { id: 'map' as Section, label: 'Map' },
          ]"
          :key="tab.id"
          type="button"
          class="technical-btn px-3 py-1.5 text-sm"
          :class="activeSection === tab.id ? 'technical-btn-primary' : ''"
          @click="activeSection = tab.id"
        >
          {{ tab.label }}
        </button>
        <div class="ml-auto flex items-center gap-2">
          <span
            class="technical-badge text-xs"
            :class="
              status === 'open'
                ? 'technical-badge-success'
                : status === 'connecting'
                  ? 'technical-badge-warning'
                  : 'technical-badge-error'
            "
          >
            {{ status }}
          </span>
          <span v-if="error" class="text-sm text-[var(--color-error)]">{{ error }}</span>
          <button type="button" class="technical-btn technical-btn-primary text-sm" @click="api.udpClients.start()">
            Start Client
          </button>
          <button
            type="button"
            :class="[buttonConfig.class, 'technical-btn text-sm']"
            :disabled="buttonConfig.disabled"
            @click="handleServerAction"
          >
            {{ buttonConfig.label }}
          </button>
          <button type="button" class="technical-btn text-sm" @click="connect">Connect</button>
          <button type="button" class="technical-btn text-sm" @click="close">Close</button>
        </div>
      </nav>
    </header>

    <main class="flex-1 min-h-0 p-4 overflow-auto">
      <div v-show="activeSection === 'log'" class="space-y-4">
        <WsLog :lines="lines" />
        <WsSendInput :disabled="!canSend" @send="handleSend" />
      </div>

      <div v-show="activeSection === 'server'" class="technical-panel p-4 max-w-xl">
        <ServerState />
        <div class="mt-4 flex gap-2">
          <button
            type="button"
            :class="[buttonConfig.class, 'technical-btn']"
            :disabled="buttonConfig.disabled"
            @click="handleServerAction"
          >
            {{ buttonConfig.label }}
          </button>
        </div>
      </div>

      <div v-show="activeSection === 'clients'" class="flex gap-4 h-full">
        <button
          type="button"
          class="technical-btn technical-btn-primary self-start"
          @click="clientsOverlay = true"
        >
          Clients öffnen
        </button>
      </div>

      <div v-show="activeSection === 'map'" class="h-full min-h-[400px]">
        <ClientMap />
      </div>
    </main>

    <Overlay
      v-model="clientsOverlay"
      title="UDP-Clients"
      width-class="w-11/12 w-[90vw]"
      max-w-class="max-w-none"
      height-class="h-11/12"
    >
      <ClientsState
        :needs-update="newClient"
        :usu-event="usuEvent"
        @done:clients-update="newClient = false"
      />
    </Overlay>
  </div>
</template>
