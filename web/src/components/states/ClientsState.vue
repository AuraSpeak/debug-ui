<template>
    <div class="h-full flex gap-4">
        <ClientList
            :clients="clients"
            :selected-client-id="clientState?.id"
            @select-client="getClientState"
            @start-server="handleStartServer"
        />
        <ClientDetails
            :client-state="clientState"
            :client-name="clientName"
            @update-state="handleStateUpdate"
        />
    </div>
</template>

<script setup lang="ts">
import { useApi } from "@/api/useApi";
import type { UDPClientState } from "@/api/types";
import { useClientState } from "@/composables/useClientState";
import ClientList from "./clients/ClientList.vue";
import ClientDetails from "./clients/ClientDetails.vue";

const props = defineProps<{
    needsUpdate: boolean;
    usuEvent: { id: number, seq: number } | null;
}>();

const emit = defineEmits<{
    (e: "done:clientStateUpdate"): void;
    (e: "done:clientsUpdate"): void;
}>();

const api = useApi();

const {
    clients,
    clientName,
    clientState,
    getClientState,
} = useClientState(
    () => props.needsUpdate,
    () => props.usuEvent,
    () => emit("done:clientStateUpdate"),
    () => emit("done:clientsUpdate")
);

function handleStartServer() {
    api.udpClients.start();
}

function handleStateUpdate(updatedState: UDPClientState) {
    clientState.value = updatedState;
}
</script>