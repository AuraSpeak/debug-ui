<template>
    <div class="space-y-2">
        <h3 class="font-semibold text-lg">Client Statistiken</h3>
        <div class="stats stats-vertical shadow w-full">
            <div class="stat">
                <div class="stat-title">Status</div>
                <div class="stat-value text-sm">{{ clientState.running ? 'Aktiv' : 'Inaktiv' }}</div>
            </div>
            <div class="stat">
                <div class="stat-title">Anzahl Datagramme</div>
                <div class="stat-value text-sm">{{ clientState.datagrams.length }}</div>
            </div>
            <div class="stat">
                <div class="stat-title">Client → Server</div>
                <div class="stat-value text-sm">
                    {{ clientState.datagrams.filter(d => d.direction === 1).length }}
                </div>
            </div>
            <div class="stat">
                <div class="stat-title">Server → Client</div>
                <div class="stat-value text-sm">
                    {{ clientState.datagrams.filter(d => d.direction === 2).length }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { UDPClientState } from "@/api/types";

defineProps<{
    clientState: UDPClientState;
}>();
</script>
