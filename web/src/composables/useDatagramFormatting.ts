/**
 * Composable für die Formatierung von Datagramm-Nachrichten
 */
export function useDatagramFormatting() {
    /**
     * Formatiert eine Datagramm-Nachricht als Hex-String
     */
    function formatDatagram(message: any): string {
        if (!message) return "";

        if (typeof message === "string") {
            const bin = atob(message);
            const bytes = Uint8Array.from(bin, (c) => c.charCodeAt(0));
            return Array.from(bytes).map(b => b.toString(16).padStart(2, "0")).join(" ").toUpperCase();
        }

        const bytes = message instanceof Uint8Array ? message : Uint8Array.from(message);
        return Array.from(bytes).map(b => b.toString(16).padStart(2, "0")).join(" ").toUpperCase();
    }

    /**
     * Formatiert eine Datagramm-Nachricht als Text-String (UTF-8)
     */
    function formatDatagramText(message: any): string {
        if (!message) return "";

        let bytes: Uint8Array;
        
        if (typeof message === "string") {
            // Base64-String vom Backend dekodieren
            try {
                const bin = atob(message);
                bytes = Uint8Array.from(bin, (c) => c.charCodeAt(0));
            } catch (e) {
                return message; // Falls kein Base64, einfach den String zurückgeben
            }
        } else if (message instanceof Uint8Array) {
            bytes = message;
        } else if (Array.isArray(message)) {
            bytes = Uint8Array.from(message);
        } else {
            return String(message);
        }

        // UTF-8 dekodieren
        try {
            const decoder = new TextDecoder('utf-8', { fatal: false });
            return decoder.decode(bytes);
        } catch (e) {
            // Falls UTF-8 Dekodierung fehlschlägt, als Hex anzeigen
            return Array.from(bytes).map(b => b.toString(16).padStart(2, "0")).join(" ");
        }
    }

    return {
        formatDatagram,
        formatDatagramText,
    };
}
