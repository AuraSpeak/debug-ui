<script setup lang="ts">
import { ref, watch, onMounted, inject, computed, type Ref } from "vue";
import { useApi } from "@/api/useApi";
import type { ClientMapData } from "@/api/types";

const SERVER_ID = 0;
const CENTER_X = 200;
const CENTER_Y = 150;
const RADIUS = 120;

const api = useApi();
const mapData = ref<ClientMapData | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

const mapRefreshTrigger = inject<Ref<number>>("mapRefreshTrigger");
const packetEventRef = inject<Ref<{ from: number; to: number; dir: number } | null>>("packetEvent");

async function fetchMap() {
  loading.value = true;
  error.value = null;
  try {
    mapData.value = await api.udpClients.getClientMap();
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : String(e);
    mapData.value = null;
  } finally {
    loading.value = false;
  }
}

onMounted(() => fetchMap());

if (mapRefreshTrigger) {
  watch(mapRefreshTrigger, () => fetchMap());
}

const pulseClientId = ref<number | null>(null);
if (packetEventRef) {
  watch(
    packetEventRef,
    (ev) => {
      if (!ev) return;
      const clientId = ev.from === SERVER_ID ? ev.to : ev.from;
      pulseClientId.value = clientId;
      packetEventRef.value = null;
      setTimeout(() => {
        pulseClientId.value = null;
      }, 400);
    },
    { deep: true }
  );
}

function clientPosition(index: number, total: number) {
  const angle = (2 * Math.PI * index) / Math.max(1, total) - Math.PI / 2;
  return {
    x: CENTER_X + RADIUS * Math.cos(angle),
    y: CENTER_Y + RADIUS * Math.sin(angle),
  };
}

const edges = computed(() => {
  const data = mapData.value;
  if (!data?.clients?.length) return [];
  return data.connections
    .filter((c) => c.toClientId === SERVER_ID)
    .map((c) => {
      const idx = data.clients.findIndex((x) => x.id === c.fromClientId);
      const pos = idx >= 0 ? clientPosition(idx, data.clients.length) : { x: CENTER_X, y: CENTER_Y };
      return { fromId: c.fromClientId, x1: pos.x, y1: pos.y, x2: CENTER_X, y2: CENTER_Y };
    });
});
</script>

<template>
  <div class="technical-panel h-full min-h-[360px] p-4 flex flex-col">
    <h2 class="text-lg font-semibold mb-2 text-[var(--color-text)]">Client-Map</h2>
    <p class="text-sm text-[var(--color-text-muted)] mb-4">
      Welcher Client mit welchem über den Server kommuniziert. Pakete werden in Echtzeit angezeigt.
    </p>

    <div v-if="loading" class="flex-1 flex items-center justify-center text-[var(--color-text-muted)]">
      Lade…
    </div>
    <div v-else-if="error" class="flex-1 flex items-center justify-center text-[var(--color-error)]">
      {{ error }}
    </div>
    <div v-else-if="mapData" class="flex-1 min-h-0 relative">
      <svg
        viewBox="0 0 400 300"
        class="w-full h-full border border-[var(--color-border)] rounded-[var(--radius)] bg-[var(--color-bg-elevated)]"
        preserveAspectRatio="xMidYMid meet"
      >
        <!-- Edges first (under nodes) -->
        <line
          v-for="e in edges"
          :key="e.fromId"
          :x1="e.x1"
          :y1="e.y1"
          :x2="e.x2"
          :y2="e.y2"
          :stroke="pulseClientId === e.fromId ? 'var(--color-accent)' : 'var(--color-border)'"
          :stroke-width="pulseClientId === e.fromId ? 2.5 : 1"
          class="transition-all duration-300"
        />

        <!-- Server node (center) -->
        <g :transform="`translate(${CENTER_X}, ${CENTER_Y})`">
          <circle r="24" fill="var(--color-bg)" stroke="var(--color-border-strong)" stroke-width="2" />
          <text text-anchor="middle" dy="0.35em" class="text-xs font-mono fill-[var(--color-text)]">Server</text>
        </g>

        <!-- Client nodes -->
        <g
          v-for="(client, i) in mapData.clients"
          :key="client.id"
          :transform="`translate(${clientPosition(i, mapData.clients.length).x}, ${clientPosition(i, mapData.clients.length).y})`"
        >
          <circle
            r="20"
            fill="var(--color-bg)"
            stroke="var(--color-accent)"
            stroke-width="1.5"
          />
          <text
            text-anchor="middle"
            dy="-1.4em"
            class="text-xs font-mono fill-[var(--color-text)]"
          >
            {{ client.name }}
          </text>
          <text
            text-anchor="middle"
            dy="0.35em"
            class="text-[10px] font-mono fill-[var(--color-text-muted)]"
          >
            #{{ client.id }}
          </text>
        </g>
      </svg>
    </div>
    <div v-else class="flex-1 flex items-center justify-center text-[var(--color-text-muted)]">
      Keine Daten.
    </div>
  </div>
</template>
