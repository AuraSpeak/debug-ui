<template>
    <dialog ref="dlg" class="p-0 bg-transparent w-full h-full max-w-none max-h-none m-0 flex justify-center items-center border-none [&:not([open])]:hidden" @close="emitClosed" @cancel.prevent="onEsc" @click="onBackdropClick">
        <div 
        ref="panel" 
        :class="[
            'bg-base-100 rounded-2xl shadow-xl overflow-hidden flex flex-col',
            widthClass,
            maxWClass,
            maxHClass,
            heightClass
        ]">
            <div v-if="title" class="px-6 pt-6 flex-shrink-0">
                <div class="flex items-start justify-between gap-4">
                    <h2 class="text-xl font-bold">{{ title }}</h2>
                    <button class="btn btn-ghost btn-sm" type="button" @click="close">✕</button>
                </div>
            </div>

            <div :class="[title ? 'px-6 pb-6 pt-4' : 'p-6', 'flex-1 min-h-0 flex flex-col']">
                <slot />
            </div>

            <div v-if="$slots.footer" class="px-6 pb-6 pt-0 flex-shrink-0">
                <slot name="footer" />
            </div>
        </div>
    </dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from "vue";

const props = defineProps({
    modelValue: { type: Boolean, default: false }, // v-model
    title: { type: String, default: "" },

    // Width control (the whole point)
    widthClass: { type: String, default: "w-[90vw]" },
    maxWClass: { type: String, default: "max-w-none" },

    // Height / scrolling
    maxHClass: { type: String, default: "max-h-[85vh] overflow-y-auto" },
    heightClass: { type: String, default: "h-full" },
    // Behavior
    closeOnBackdrop: { type: Boolean, default: true },
    closeOnEsc: { type: Boolean, default: true },
});

const emit = defineEmits(["update:modelValue"]);

const dlg = ref<HTMLDialogElement | null>(null);
const panel = ref<HTMLDivElement | null>(null);

function open() {
    if (dlg.value && !dlg.value.open) dlg.value.showModal();
}

function close() {
    // close() triggers the native 'close' event; we sync v-model there too.
    if (dlg.value && dlg.value.open) dlg.value.close();
}

function emitClosed() {
    // keep parent state in sync (covers ESC, backdrop click, X button, etc.)
    emit("update:modelValue", false);
}

function onEsc() {
    if (!props.closeOnEsc) return;
    close();
}

function onBackdropClick(e: MouseEvent) {
    if (!props.closeOnBackdrop) return;

    // If click is outside panel -> close
    const r = panel.value?.getBoundingClientRect() as DOMRect;
    if (!r) return;

    const inside =
        e.clientX >= r.left &&
        e.clientX <= r.right &&
        e.clientY >= r.top &&
        e.clientY <= r.bottom;

    if (!inside) close();
}

// Sync v-model -> dialog open/close
watch(
    () => props.modelValue,
    (v) => (v ? open() : close()),
    { immediate: true }
);

// Safety: if component unmounts while open, close dialog to avoid stuck backdrop
onBeforeUnmount(() => {
    if (dlg.value?.open) dlg.value.close();
});
</script>

<style scoped>
/* Backdrop for native <dialog> - using Tailwind bg-black/60 */
dialog::backdrop {
    background: rgba(0, 0, 0, 0.6);
}
</style>