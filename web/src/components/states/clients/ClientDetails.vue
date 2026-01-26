<template>
    <div class="flex-1 border border-base-300 rounded-lg p-4 bg-base-200 h-full overflow-y-auto">
        <div v-if="clientState" class="space-y-4">
            <ClientHeader
                :client-name="clientName"
                :client-id="clientState.id"
                :is-running="clientState.running"
            />
            <ClientStats :client-state="clientState" />
            <SendDatagramForm
                :client-id="clientState.id"
                @sent="handleDatagramSent"
            />
            <DatagramTable :datagrams="clientState.datagrams" />
            <ClientMermaid :cname="clientName ?? ''" />
        </div>
        <div v-else class="flex items-center justify-center h-full">
            <div class="text-center">
                <p class="text-xl font-semibold text-base-content/60">Kein Client ausgewählt</p>
                <p class="text-sm text-base-content/40 mt-2">Wählen Sie einen Client aus der Liste aus</p>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { UDPClientState } from "@/api/types";
import ClientHeader from "./ClientHeader.vue";
import ClientStats from "./ClientStats.vue";
import SendDatagramForm from "./SendDatagramForm.vue";
import DatagramTable from "./DatagramTable.vue";
import ClientMermaid from "./ClientMermaid.vue";

const props = defineProps<{
    clientState: UDPClientState | undefined;
    clientName: string | null;
}>();

const emit = defineEmits<{
    (e: "update-state", updatedState: UDPClientState): void;
}>();

function handleDatagramSent(updatedState: UDPClientState) {
    emit("update-state", updatedState);
}
</script>
