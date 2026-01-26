<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  disabled?: boolean;
}>();

const emit = defineEmits<{
  (e: "send", value: string): void;
}>();

const msg = ref("");

function submit() {
  const t = msg.value.trim();
  if (!t) return;
  emit("send", t);
  msg.value = "";
}
</script>

<template>
  <div class="join w-full">
    <input
      class="input input-bordered join-item w-full"
      placeholder="Nachricht…"
      v-model="msg"
      @keyup.enter="submit"
      :disabled="disabled"
    />
    <button class="btn btn-primary join-item" @click="submit" :disabled="disabled">
      Senden
    </button>
  </div>
</template>
