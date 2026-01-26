<template>
    <div class="w-1/4 border border-base-300 rounded-lg p-4 bg-base-200 flex flex-col h-full">
        <h2 class="text-xl font-bold mb-4">Clients</h2>
        <div class="flex flex-col gap-2 overflow-y-auto flex-1 min-h-0">
            <button
                v-for="client in clients"
                :key="client.id"
                class="btn w-full justify-start"
                :class="selectedClientId === client.id ? 'btn-primary' : 'btn-outline'"
                @click="$emit('select-client', client.id)"
            >
                <div class="flex items-center gap-2 w-full">
                    <div class="flex-1 text-left">{{ client.name }}</div>
                    <div class="badge badge-sm badge-ghost">#{{ client.id }}</div>
                </div>
            </button>
        </div>
        <button class="btn btn-primary w-full mt-4" @click="$emit('start-server')">
            Start Client
        </button>
    </div>
</template>

<script setup lang="ts">
import type { UDPClient } from "@/api/types";

defineProps<{
    clients: UDPClient[];
    selectedClientId?: number;
}>();

defineEmits<{
    (e: "select-client", id: number): void;
    (e: "start-server"): void;
}>();
</script>
