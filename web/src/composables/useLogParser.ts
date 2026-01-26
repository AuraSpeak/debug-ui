import { computed, ref, type Ref } from "vue";
import type { LogEntry } from "@/api/types";

export interface ParsedLogLine {
  type: "log" | "text";
  log?: LogEntry;
  text?: string;
  index: number;
}

/**
 * Composable zum Parsen von WebSocket-Nachrichten in Log-Einträge und Text-Nachrichten
 */
export function useLogParser(lines: Ref<string[]>) {
  const parsedLines = computed<ParsedLogLine[]>(() => {
    return lines.value.map((line, index) => {
      // Trim Whitespace (logrus JSONFormatter fügt oft ein Newline hinzu)
      const trimmedLine = line.trim();
      
      // Überspringe leere Zeilen
      if (!trimmedLine) {
        return {
          type: "text",
          text: line,
          index,
        };
      }
      
      // Versuche JSON zu parsen
      try {
        const parsed = JSON.parse(trimmedLine);
        // Prüfe ob es ein Log-Eintrag ist (hat level, msg)
        if (
          typeof parsed === "object" &&
          parsed !== null &&
          typeof parsed.level === "string" &&
          typeof parsed.msg === "string"
        ) {
          return {
            type: "log",
            log: parsed as LogEntry,
            index,
          };
        }
      } catch {
        // Kein gültiges JSON, behandle als Text
      }

      return {
        type: "text",
        text: line,
        index,
      };
    });
  });

  const logEntries = computed<LogEntry[]>(() => {
    return parsedLines.value
      .filter((line) => line.type === "log" && line.log)
      .map((line) => line.log!);
  });

  const textLines = computed<string[]>(() => {
    return parsedLines.value
      .filter((line) => line.type === "text" && line.text)
      .map((line) => line.text!);
  });

  return {
    parsedLines,
    logEntries,
    textLines,
  };
}
