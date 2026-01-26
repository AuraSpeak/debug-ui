<template>
    <div tabindex="0" class="collapse bg-base-100">
        <h3 class="font-semibold text-lg collapse-title">Datagramm senden</h3>
        <div class="card bg-base-100 shadow collapse-content">
            <div class="card-body p-4">
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Format</span>
                    </label>
                    <div class="flex gap-4">
                        <label class="label cursor-pointer gap-2">
                            <input
                                type="radio"
                                name="format"
                                class="radio radio-primary"
                                value="text"
                                v-model="sendFormat"
                            />
                            <span class="label-text">Text</span>
                        </label>
                        <label class="label cursor-pointer gap-2">
                            <input
                                type="radio"
                                name="format"
                                class="radio radio-primary"
                                value="hex"
                                v-model="sendFormat"
                            />
                            <span class="label-text">Hex</span>
                        </label>
                    </div>
                </div>
                <div class="form-control">
                    <label class="label">
                        <span class="label-text">Nachricht:</span>
                    </label>
                    <br>
                    <textarea
                        v-model="sendMessage"
                        class="textarea textarea-bordered w-full resize-none"
                        :placeholder="sendFormat === 'hex' ? '48 65 6C 6C 6F oder 48656C6C6F' : 'Geben Sie Ihre Nachricht ein...'"
                        rows="3"
                    ></textarea>
                </div>
                <div class="form-control mt-2">
                    <button
                        class="btn btn-primary"
                        :disabled="!sendMessage.trim() || sending"
                        @click="handleSend"
                        :key="`send-btn-${sending}`"
                    >
                        <span v-if="sending" class="loading loading-spinner loading-sm" key="spinner"></span>
                        <span v-else key="send-text">Senden</span>
                    </button>
                </div>
                <div v-if="sendError" class="alert alert-error mt-2">
                    <span>{{ sendError }}</span>
                </div>
                <div v-if="sendSuccess" class="alert alert-success mt-2">
                    <span>Datagramm erfolgreich gesendet!</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useSendDatagram } from "@/composables/useSendDatagram";

const props = defineProps<{
    clientId: number;
}>();

const emit = defineEmits<{
    (e: "sent", updatedState: any): void;
}>();

const {
    sendFormat,
    sendMessage,
    sending,
    sendError,
    sendSuccess,
    sendDatagram,
} = useSendDatagram();

async function handleSend() {
    await sendDatagram(props.clientId, (updatedState) => {
        emit("sent", updatedState);
    });
}
</script>
