<template>
    <div class="space-y-2" v-if="datagrams.length > 0">
        <h3 class="font-semibold text-lg">Datagramme</h3>
        <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
                <thead>
                    <tr>
                        <th>Richtung</th>
                        <th>Nachricht (Hex)</th>
                        <th>Nachricht (String)</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(datagram, index) in datagrams" :key="index">
                        <td class="w-1/4">
                            <div class="badge"
                                :class="datagram.direction === 1 ? 'badge-info' : 'badge-warning'">
                                {{ datagram.direction === 1 ? 'Client → Server' : 'Server → Client' }}
                            </div>
                        </td>
                        <td>
                            <code class="text-xs">
                                {{ formatDatagram(datagram.message) }}
                            </code>
                        </td>
                        <td>
                            <code class="text-xs">
                                {{ formatDatagramText(datagram.message) }}
                            </code>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    <div v-else class="alert alert-info">
        <span>Keine Datagramme vorhanden</span>
    </div>
</template>

<script setup lang="ts">
import type { Datagram } from "@/api/types";
import { useDatagramFormatting } from "@/composables/useDatagramFormatting";

defineProps<{
    datagrams: Datagram[];
}>();

const { formatDatagram, formatDatagramText } = useDatagramFormatting();
</script>
