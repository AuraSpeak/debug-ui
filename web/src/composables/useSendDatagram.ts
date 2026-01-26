import { ref } from "vue";
import { useApi } from "@/api/useApi";
import type { UDPClientState } from "@/api/types";

/**
 * Composable für das Senden von Datagrammen
 */
export function useSendDatagram() {
    const api = useApi();
    
    const sendFormat = ref<"hex" | "text">("text");
    const sendMessage = ref<string>("");
    const sending = ref<boolean>(false);
    const sendError = ref<string | null>(null);
    const sendSuccess = ref<boolean>(false);

    /**
     * Sendet ein Datagramm an einen Client
     */
    async function sendDatagram(
        clientId: number,
        onSuccess?: (updatedState: UDPClientState) => void
    ): Promise<void> {
        if (!sendMessage.value.trim()) {
            return;
        }

        sending.value = true;
        sendError.value = null;
        sendSuccess.value = false;

        try {
            await api.udpClients.sendDatagram({
                id: clientId,
                message: sendMessage.value.trim(),
                format: sendFormat.value,
            });

            sendSuccess.value = true;
            sendMessage.value = "";

            // Update client state after sending
            const updatedState = await api.udpClients.getStateById(clientId);
            onSuccess?.(updatedState);

            // Clear success message after 3 seconds
            setTimeout(() => {
                sendSuccess.value = false;
            }, 3000);
        } catch (error: any) {
            sendError.value = error.message || "Fehler beim Senden des Datagramms";
            console.error("Error sending datagram:", error);
        } finally {
            sending.value = false;
        }
    }

    return {
        sendFormat,
        sendMessage,
        sending,
        sendError,
        sendSuccess,
        sendDatagram,
    };
}
