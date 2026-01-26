package util

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"

	"github.com/auraspeak/server/pkg/tracer"
)

// ---- Mermaid builder ----

type MermaidBuilder struct {
	sb    strings.Builder
	nodes map[string]string // label -> nodeID
}

func NewMermaidBuilder() *MermaidBuilder {
	b := &MermaidBuilder{
		nodes: make(map[string]string, 64),
	}
	b.sb.WriteString("flowchart TD\n")
	return b
}

func (b *MermaidBuilder) Node(label string) string {
	if id, ok := b.nodes[label]; ok {
		return id
	}
	id := "n" + hashID(label) // stable, mermaid-safe id
	b.nodes[label] = id

	// Define node line: n123["Label"]
	b.sb.WriteString("  ")
	b.sb.WriteString(id)
	b.sb.WriteString("[\"")
	b.sb.WriteString(escapeMermaidLabel(label))
	b.sb.WriteString("\"]\n")

	return id
}

func (b *MermaidBuilder) Edge(fromID, toID, label string) {
	// Edge line: a -->|label| b
	b.sb.WriteString("  ")
	b.sb.WriteString(fromID)
	b.sb.WriteString(" -->")

	if label != "" {
		b.sb.WriteString("|\"")
		b.sb.WriteString(escapeMermaidLabel(label))
		b.sb.WriteString("\"|")
	}

	b.sb.WriteString(" ")
	b.sb.WriteString(toID)
	b.sb.WriteString("\n")
}

func (b *MermaidBuilder) String() string { return b.sb.String() }

// ---- Conversion: []TraceEvent -> Mermaid ----

func BuildMermaidFromTraces(events []tracer.TraceEvent) string {
	// Optional: stable ordering (nice for diffs / repeatability)
	sort.Slice(events, func(i, j int) bool {
		if events[i].TS.Equal(events[j].TS) {
			return i < j
		}
		return events[i].TS.Before(events[j].TS)
	})

	mb := NewMermaidBuilder()

	for _, ev := range events {
		localID := mb.Node("local: " + ev.Local)
		remoteID := mb.Node("remote: " + ev.Remote)

		edgeLabel := fmt.Sprintf(
			"%s len=%d cid=%d",
			ev.TS.Format("15:04:05.000"),
			ev.Len,
			ev.ClientID,
		)

		switch ev.Dir {
		case tracer.TraceIn:
			mb.Edge(remoteID, localID, "IN "+edgeLabel)
		case tracer.TraceOut:
			mb.Edge(localID, remoteID, "OUT "+edgeLabel)
		default:
			mb.Edge(remoteID, localID, string(ev.Dir)+" "+edgeLabel)
		}
	}

	return mb.String()
}

// Sequence diagram: Client -> Server / Server -> Client
func BuildSequenceDiagramFromTraces(events []tracer.TraceEvent) string {
	// stabil und zeitlich korrekt
	sort.Slice(events, func(i, j int) bool {
		if events[i].TS.Equal(events[j].TS) {
			return i < j
		}
		return events[i].TS.Before(events[j].TS)
	})

	var sb strings.Builder
	sb.WriteString("sequenceDiagram\n")

	// Server participant (Local ist bei dir der UDP-Server-Socket)
	serverLabel := "Server"
	if len(events) > 0 && events[0].Local != "" {
		serverLabel = "Server " + events[0].Local
	}
	sb.WriteString("  participant S as \"")
	sb.WriteString(escapeMermaidLabel(serverLabel))
	sb.WriteString("\"\n")

	// Falls Remote sich ändert: pro Remote einen eigenen Client-Participant
	clientPID := map[string]string{} // remote -> pid

	getClientPID := func(remote string, cid int) string {
		if remote == "" {
			remote = "unknown"
		}
		if pid, ok := clientPID[remote]; ok {
			return pid
		}
		pid := "C" + hashID(remote)
		clientPID[remote] = pid

		// Participant definieren
		lbl := fmt.Sprintf("Client cid=%d\\n%s", cid, remote)
		sb.WriteString("  participant ")
		sb.WriteString(pid)
		sb.WriteString(" as \"")
		sb.WriteString(escapeMermaidLabel(lbl))
		sb.WriteString("\"\n")

		return pid
	}

	for _, ev := range events {
		c := getClientPID(ev.Remote, ev.ClientID)

		label := fmt.Sprintf("%s len=%d cid=%d",
			ev.TS.Format("15:04:05.000"),
			ev.Len,
			ev.ClientID,
		)

		switch ev.Dir {
		case tracer.TraceIn:
			// remote -> local  (Client sendet an Server)
			sb.WriteString("  ")
			sb.WriteString(c)
			sb.WriteString("->>S: ")
			sb.WriteString(escapeMermaidLabel("SEND " + label))
			sb.WriteString("\n")

		case tracer.TraceOut:
			// local -> remote  (Server sendet an Client)
			sb.WriteString("  S->>")
			sb.WriteString(c)
			sb.WriteString(": ")
			sb.WriteString(escapeMermaidLabel("SEND " + label))
			sb.WriteString("\n")

		default:
			sb.WriteString("  Note over S,")
			sb.WriteString(c)
			sb.WriteString(": ")
			sb.WriteString(escapeMermaidLabel("dir=" + string(ev.Dir) + " " + label))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ---- Helpers ----

func hashID(s string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum32())
}

// Escapes characters that commonly break Mermaid labels when put inside quotes.
// Keep it conservative; Mermaid parsing can be picky.
func escapeMermaidLabel(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}
