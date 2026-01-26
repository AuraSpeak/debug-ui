<script setup lang="ts">
import { computed, ref, toRef } from "vue";
import { useLogParser } from "@/composables/useLogParser";
import LogDetailModal from "./LogDetailModal.vue";
import type { LogEntry } from "@/api/types";

const props = defineProps<{
  lines: string[];
}>();

const { parsedLines } = useLogParser(toRef(props, "lines"));

const selectedLog = ref<LogEntry | null>(null);
const searchQuery = ref("");

const filteredLines = computed(() => {
  if (!searchQuery.value.trim()) {
    return parsedLines.value;
  }

  const query = searchQuery.value.toLowerCase().trim();
  
  return parsedLines.value.filter((line) => {
    if (line.type === "log" && line.log) {
      // Suche in msg, level, caller und allen anderen Feldern
      const msg = (line.log.msg || "").toLowerCase();
      const level = (line.log.level || "").toLowerCase();
      const caller = (line.log.caller || "").toLowerCase();
      
      // Suche auch in allen anderen Feldern des Log-Eintrags
      const allFields = Object.values(line.log)
        .map((val) => String(val).toLowerCase())
        .join(" ");
      
      return (
        msg.includes(query) ||
        level.includes(query) ||
        caller.includes(query) ||
        allFields.includes(query)
      );
    } else if (line.type === "text" && line.text) {
      return line.text.toLowerCase().includes(query);
    }
    return false;
  });
});

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

function getCallerBadgeClass(caller: string | undefined): string {
  if (!caller) return "badge-ghost";
  const callerLower = caller.toLowerCase();
  if (callerLower === "web") {
    return "badge-primary";
  }
  if (callerLower === "server") {
    return "badge-secondary";
  }
  if (callerLower === "protocol") {
    return "badge-accent";
  }
  if (callerLower === "client") {
    return "badge-info";
  }
  return "badge-ghost";
}

function getCallerLabel(caller: string | undefined): string {
  if (!caller) return "?";
  return caller.toUpperCase();
}

function openLogDetail(entry: LogEntry) {
  selectedLog.value = entry;
}

function closeLogDetail() {
  selectedLog.value = null;
}
</script>

<template>
  <div class="bg-base-200 rounded-lg h-72 flex flex-col p-3">
    <!-- Suchleiste -->
    <div class="mb-3">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Logs durchsuchen..."
        class="input input-bordered input-sm w-full"
      />
    </div>

    <!-- Legende -->
    <div class="mb-3 pb-2 border-b border-base-300">
      <div class="text-xs font-semibold text-base-content/70 mb-1">Caller:</div>
      <div class="flex flex-wrap gap-2">
        <div class="flex items-center gap-1">
          <span class="badge badge-sm badge-primary">WEB</span>
          <span class="text-xs text-base-content/60">Web Server</span>
        </div>
        <div class="flex items-center gap-1">
          <span class="badge badge-sm badge-secondary">SERVER</span>
          <span class="text-xs text-base-content/60">UDP Server</span>
        </div>
        <div class="flex items-center gap-1">
          <span class="badge badge-sm badge-accent">PROTOCOL</span>
          <span class="text-xs text-base-content/60">Protocol</span>
        </div>
        <div class="flex items-center gap-1">
          <span class="badge badge-sm badge-info">CLIENT</span>
          <span class="text-xs text-base-content/60">UDP Client</span>
        </div>
      </div>
    </div>

    <!-- Log-Einträge -->
    <div class="flex-1 overflow-auto space-y-1">
      <div
        v-if="filteredLines.length === 0"
        class="text-center text-sm text-base-content/50 py-4"
      >
        Keine Logs gefunden
      </div>
      <div
        v-for="line in filteredLines"
        :key="line.index"
        class="flex items-start gap-2"
      >
        <div
          v-if="line.type === 'log' && line.log"
          class="flex-1 cursor-pointer hover:bg-base-300 rounded px-2 py-1 transition-colors border border-base-300"
          @click="openLogDetail(line.log)"
        >
          <div class="flex items-center gap-2">
            <span
              class="badge badge-sm"
              :class="getCallerBadgeClass(line.log.caller)"
              :title="`Caller: ${line.log.caller || 'unknown'}`"
            >
              {{ getCallerLabel(line.log.caller) }}
            </span>
            <span
              class="badge badge-sm"
              :class="getLevelBadgeClass(line.log.level)"
            >
              {{ line.log.level }}
            </span>
            <span class="text-sm">{{ line.log.msg }}</span>
          </div>
        </div>
        <div
          v-else-if="line.type === 'text' && line.text"
          class="flex-1 text-sm whitespace-pre-wrap border border-base-300 rounded px-2 py-1"
        >
          {{ line.text }}
        </div>
      </div>
    </div>
    
    <LogDetailModal :entry="selectedLog" @close="closeLogDetail" />
  </div>
</template>
