<template>
    <div class="flex items-center justify-between">
        <div>
            <h2 class="text-2xl font-bold">{{ clientName || 'Unbekannter Client' }}</h2>
            <p class="text-sm text-base-content/60">ID: {{ clientId }}</p>
        </div>
        <div class="badge badge-lg"
            :class="isRunning ? 'badge-success' : 'badge-error'">
            {{ isRunning ? 'Läuft' : 'Gestoppt' }}
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    clientName: string | null;
    clientId: number;
    isRunning: boolean;
}>();
</script>
