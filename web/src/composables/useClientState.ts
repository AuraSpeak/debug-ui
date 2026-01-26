import { ref, watch, onMounted } from "vue";
import { useApi } from "@/api/useApi";
import type { UDPClient, UDPClientState } from "@/api/types";

/**
 * Composable für das State-Management von UDP-Clients
 */
export function useClientState(
    needsUpdate: () => boolean,
    usuEvent: () => { id: number, seq: number } | null,
    onClientStateUpdate?: () => void,
    onClientsUpdate?: () => void
) {
    const api = useApi();
    const clients = ref<UDPClient[]>([]);
    const clientName = ref<string | null>(null);
    const clientState = ref<UDPClientState>();
    const usuSeq = ref<number | null>(null);

    /**
     * Lädt alle Clients vom Server
     */
    async function refreshClients(): Promise<void> {
        const clientData = await api.udpClients.getAll();
        clients.value = clientData.udpClients;
        console.log("clients updated", clients.value);
        onClientsUpdate?.();
    }

    /**
     * Lädt den State eines spezifischen Clients
     */
    async function getClientState(id: number): Promise<void> {
        clientName.value = clients.value.find(c => c.id === id)?.name || null;
        const clientStateData = await api.udpClients.getStateById(id);
        console.log(clientStateData);
        clientState.value = clientStateData;
    }

    /**
     * Behandelt Update-Events vom Server
     */
    async function handleUsUEvent(e: { id: number, seq: number }): Promise<void> {
        console.log("handleUsUEvent", e);
        if (usuSeq.value === e.seq) return;
        usuSeq.value = e.seq;
        if (e.id === clientState.value?.id) {
            console.log("handleUsUEvent", "updating client state");
            const StateData = await api.udpClients.getStateById(e.id);
            clientState.value = StateData;
            onClientStateUpdate?.();
        }
    }

    // Watcher für Update-Events
    watch(
        usuEvent,
        (e) => {
            if (!e) return;
            handleUsUEvent(e);
        },
        { flush: "post" }
    );

    // Watcher für Client-Updates
    watch(needsUpdate, async () => {
        await refreshClients();
    });

    // Initial load
    onMounted(async () => {
        await refreshClients();
    });

    return {
        clients,
        clientName,
        clientState,
        refreshClients,
        getClientState,
        handleUsUEvent,
    };
}
