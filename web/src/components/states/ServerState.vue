<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue';
import { useServerStore } from '@/stores/UDPServer';

const serverStore = useServerStore();
let timer: number | undefined;

onUnmounted(() => {
    if (timer) {
        clearInterval(timer);
    }
});

const lastUpdatedFormatted = computed(() => {
    if (!serverStore.lastUpdated) return 'Never updated';
    const date = new Date(serverStore.lastUpdated);
    return date.toLocaleTimeString('de-DE');
});
</script>

<template>
    <div class="space-y-4">
        <div class="flex items-center gap-2">
            <span class="badge"
                :class="serverStore.ServerState?.isAlive ? 'badge-success' : serverStore.ServerState?.shouldStop ? 'badge-warning' : 'badge-ghost'">
                {{ serverStore.ServerState?.isAlive ? 'Running' : serverStore.ServerState?.shouldStop ? 'Stopping' : 'Stopped' }}
            </span>
            <span v-if="serverStore.error" class="text-error text-sm">{{ serverStore.error }}</span>
            <span v-if="serverStore.loading" class="loading loading-spinner loading-xs"></span>
        </div>

        <div class="space-y-2 flex flex-col">
            <h3 class="font-semibold text-lg">Server Statistiken</h3>
            <div class="stats stats-vertical shadow">
                <div class="stat w-full">
                    <div class="stat-title">Status</div>
                    <div class="stat-value text-sm">{{ serverStore.status }}</div>
                </div>
                <div class="stat">
                    <div class="stat-title">Server läuft</div>
                    <div class="stat-value text-sm">{{ serverStore.ServerState?.isAlive ? 'Ja' : 'Nein' }}</div>
                </div>
                <div class="stat">
                    <div class="stat-title">Wird gestoppt</div>
                    <div class="stat-value text-sm">{{ serverStore.ServerState?.shouldStop ? 'Ja' : 'Nein' }}</div>
                </div>
                <div class="stat">
                    <div class="stat-title">Letzte Aktualisierung</div>
                    <div class="stat-value text-sm">{{ lastUpdatedFormatted }}</div>
                </div>
            </div>
        </div>
    </div>
</template>