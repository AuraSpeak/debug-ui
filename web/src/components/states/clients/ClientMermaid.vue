<template>
  <h3>Mermaid Diagrams</h3>
  <button class="btn btn-primary" @click="getWhole">Whole Lifetime</button>

  <div tabindex="0" class="collapse bg-base-100">
    <h3 class="collapse-title" v-if="heading !== ''">{{ heading }}</h3>

    <div class="card shadow collapse-content" style="background-color: lightblue;">
      <div v-if="error" class="text-error">{{ error }}</div>
      <div v-else-if="loading">Loading…</div>
      <div v-else v-html="svg"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useApi } from "@/api/useApi";
import mermaid from "mermaid";

const props = defineProps<{ cname: string }>();

const api = useApi();
const heading = ref("");
const svg = ref("");
const loading = ref(false);
const error = ref<string | null>(null);

// Funktion zum Abrufen der DaisyUI-Farben aus CSS-Variablen
function getDaisyUIColors() {
  const root = document.documentElement;
  const computedStyle = getComputedStyle(root);
  
  // DaisyUI verwendet CSS-Variablen wie --p, --s, --a, --b1, --b2, --b3, --bc, etc.
  // Fallback zu Standardwerten, falls Variablen nicht verfügbar sind
  const getColor = (varName: string, fallback: string) => {
    const value = computedStyle.getPropertyValue(varName).trim();
    return value || fallback;
  };

  // Versuche DaisyUI-Variablen zu lesen, mit Fallbacks
  const primary = getColor("--p", "#3b82f6");
  const base100 = getColor("--b1", "#ffffff");
  const base200 = getColor("--b2", "#f3f4f6");
  const base300 = getColor("--b3", "#e5e7eb");
  const baseContent = getColor("--bc", "#1f2937");
  const accent = getColor("--a", "#6366f1");
  const secondary = getColor("--s", "#8b5cf6");

  return {
    primary,
    base100,
    base200,
    base300,
    baseContent,
    accent,
    secondary,
  };
}

// Mermaid-Theme an DaisyUI anpassen
function initializeMermaid() {
  const colors = getDaisyUIColors();
  
  mermaid.initialize({
    startOnLoad: false,
    securityLevel: "strict",
    theme: "base",
    themeVariables: {
      // Hauptfarben
      primaryColor: colors.primary,
      primaryTextColor: colors.baseContent,
      primaryBorderColor: colors.base300,
      secondaryColor: colors.secondary,
      tertiaryColor: colors.accent,
      
      // Hintergrundfarben - explizite helle Farbe für bessere Sichtbarkeit
      mainBkgColor: "#ffffff",  // Weiß - hellster Hintergrund
      secondBkgColor: "#f5f5f5",  // Sehr helles Grau
      tertiaryBkgColor: "#e5e5e5",  // Helles Grau
      
      // Textfarben
      textColor: colors.baseContent,
      lineColor: colors.baseContent,
      titleColor: colors.baseContent,
      
      // Edge/Line Farben - Labels zwischen Nodes besser sichtbar machen
      edgeLabelBackground: "#ffffff",  // Weiß für Labels
      edgeLabelTextColor: colors.baseContent,
      defaultLinkColor: colors.baseContent,
      clusterBkg: colors.base200,
      clusterBorder: colors.base300,
      
      // Spezifische Diagramm-Farben (für verschiedene Diagrammtypen)
      cScale0: colors.primary,
      cScale1: colors.secondary,
      cScale2: colors.accent,
    },
  });
}

onMounted(() => {
  initializeMermaid();
});

// If your backend sometimes includes ```mermaid fences, strip them.
function normalizeMermaid(src: string) {
  return src
    .replace(/^```mermaid\s*/i, "")
    .replace(/```$/i, "")
    .trim();
}

async function getWhole() {
  loading.value = true;
  error.value = null;

  try {
    // Farben vor jedem Rendern aktualisieren
    initializeMermaid();
    
    const res = await api.trace.getAll(props.cname);

    // IMPORTANT: your Go struct tag is json:"diagrams"
    const source = normalizeMermaid(res.diagram);
    heading.value = res.heading ?? "";

    const id = "mmd-" + Math.random().toString(16).slice(2);
    const out = await mermaid.render(id, source);

    svg.value = out.svg;
  } catch (e: any) {
    error.value = e?.message ?? String(e);
    svg.value = "";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
/* Stelle sicher, dass Mermaid-Diagramme einen hellen Hintergrund haben */
:deep(.mermaid svg) {
  background-color: #ffffff !important; /* Weiß - explizite helle Farbe */
}

/* Zusätzlich: Stelle sicher, dass der Diagramm-Hintergrund selbst auch hell ist */
:deep(.mermaid svg .background) {
  fill: #ffffff !important;
}

/* Für Sequence-Diagramme: Hintergrund zwischen den Lifelines */
:deep(.mermaid svg rect[fill*="rgb"]),
:deep(.mermaid svg .background) {
  fill: #ffffff !important;
}

/* Alle Hintergrund-Elemente im SVG hell machen */
:deep(.mermaid svg rect:not(.node rect):not(.edgeLabel .label-box)) {
  fill: #ffffff !important;
}

:deep(.mermaid .node rect),
:deep(.mermaid .node circle),
:deep(.mermaid .node ellipse),
:deep(.mermaid .node polygon) {
  fill: hsl(var(--b2, 0 0% 96%));
  stroke: hsl(var(--b3, 0 0% 90%));
}

:deep(.mermaid .edgePath .path) {
  stroke: hsl(var(--bc, 0 0% 20%));
}

/* Edge-Labels (Text zwischen Nodes) besser sichtbar machen - hellerer Hintergrund */
:deep(.mermaid .edgeLabel) {
  background-color: hsl(var(--b1, 1 0% 100%)) !important;
  padding: 2px 4px;
  border-radius: 4px;
}

:deep(.mermaid .edgeLabel text) {
  fill: hsl(var(--bc, 0 0% 20%)) !important;
  font-weight: 500;
  font-size: 12px;
}

:deep(.mermaid .edgeLabel .label-box) {
  fill: hsl(var(--b1, 1 0% 100%)) !important;
  stroke: hsl(var(--b3, 0 0% 90%)) !important;
  stroke-width: 1;
}

/* Speziell für Sequence-Diagramme */
:deep(.mermaid .messageText0),
:deep(.mermaid .messageText1),
:deep(.mermaid .messageText2),
:deep(.mermaid .messageText3) {
  fill: hsl(var(--bc, 0 0% 20%)) !important;
  font-weight: 500;
}

:deep(.mermaid .messageLine0),
:deep(.mermaid .messageLine1),
:deep(.mermaid .messageLine2),
:deep(.mermaid .messageLine3) {
  stroke: hsl(var(--bc, 0 0% 20%)) !important;
  stroke-width: 1.5;
}

:deep(.mermaid .label text) {
  fill: hsl(var(--bc, 0 0% 20%));
}
</style>
