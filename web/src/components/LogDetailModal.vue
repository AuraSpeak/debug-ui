<script setup lang="ts">
import { computed, watch, ref } from "vue";
import type { LogEntry } from "@/api/types";

const props = defineProps<{
  entry: LogEntry | null;
}>();

const emit = defineEmits<{
  close: [];
}>();

const dialogRef = ref<HTMLDialogElement | null>(null);

const isOpen = computed(() => props.entry !== null);

watch(isOpen, (open) => {
  if (open && dialogRef.value) {
    dialogRef.value.showModal();
  } else if (!open && dialogRef.value) {
    dialogRef.value.close();
  }
});

function handleClose() {
  emit("close");
}

function formatValue(value: any): string {
  if (value === null) return "null";
  if (value === undefined) return "undefined";
  if (typeof value === "object") {
    return JSON.stringify(value, null, 2);
  }
  return String(value);
}

function getLevelBadgeClass(level: string): string {
  const levelLower = level.toLowerCase();
  if (levelLower === "error" || levelLower === "fatal") {
    return "badge-error";
  }
  if (levelLower === "warn" || levelLower === "warning") {
    return "badge-warning";
  }
  if (levelLower === "info") {
    return "badge-info";
  }
  if (levelLower === "debug" || levelLower === "trace") {
    return "badge-ghost";
  }
  return "badge";
}

const sortedFields = computed(() => {
  if (!props.entry) return [];
  
  // Standardfelder zuerst
  const standardFields = ["level", "msg", "time"];
  const otherFields: string[] = [];
  const entry = props.entry;
  
  for (const key in entry) {
    if (!standardFields.includes(key)) {
      otherFields.push(key);
    }
  }
  
  return [
    ...standardFields.filter((key) => key in entry),
    ...otherFields.sort(),
  ];
});
</script>

<template>
  <dialog ref="dialogRef" class="modal" @close="handleClose">
    <div class="modal-box max-w-2xl">
      <h3 class="font-bold text-lg mb-4">Log Details</h3>
      
      <div v-if="entry" class="space-y-3">
        <div
          v-for="key in sortedFields"
          :key="key"
          class="flex flex-col gap-1"
        >
          <div class="text-sm font-semibold text-base-content/70">
            {{ key }}
          </div>
          <div
            v-if="key === 'level'"
            class="badge"
            :class="getLevelBadgeClass(entry[key])"
          >
            {{ entry[key] }}
          </div>
          <pre
            v-else
            class="bg-base-200 p-2 rounded text-sm overflow-x-auto whitespace-pre-wrap"
          >
{{ formatValue(entry[key]) }}
          </pre>
        </div>
      </div>

      <div class="modal-action">
        <form method="dialog">
          <button class="btn" @click="handleClose">Close</button>
        </form>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
</template>
