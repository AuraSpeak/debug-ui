package util

import (
	"strings"
	"testing"
	"time"

	"github.com/auraspeak/server/pkg/tracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMermaidBuilder(t *testing.T) {
	mb := NewMermaidBuilder()
	require.NotNil(t, mb)
	assert.True(t, strings.HasPrefix(mb.String(), "flowchart TD"), "Should start with flowchart TD")
}

func TestMermaidBuilder_Node(t *testing.T) {
	mb := NewMermaidBuilder()

	// Add a node
	nodeID := mb.Node("TestLabel")
	require.NotEmpty(t, nodeID)
	assert.True(t, strings.HasPrefix(nodeID, "n"), "Node ID should start with 'n'")

	// Adding same label should return same ID
	nodeID2 := mb.Node("TestLabel")
	assert.Equal(t, nodeID, nodeID2, "Same label should return same node ID")

	// Adding different label should return different ID
	nodeID3 := mb.Node("DifferentLabel")
	assert.NotEqual(t, nodeID, nodeID3, "Different labels should return different node IDs")
}

func TestMermaidBuilder_Edge(t *testing.T) {
	mb := NewMermaidBuilder()

	fromID := mb.Node("From")
	toID := mb.Node("To")

	mb.Edge(fromID, toID, "TestLabel")
	result := mb.String()

	assert.Contains(t, result, fromID, "Result should contain from node ID")
	assert.Contains(t, result, toID, "Result should contain to node ID")
	assert.Contains(t, result, "TestLabel", "Result should contain edge label")
}

func TestMermaidBuilder_Edge_NoLabel(t *testing.T) {
	mb := NewMermaidBuilder()

	fromID := mb.Node("From")
	toID := mb.Node("To")

	mb.Edge(fromID, toID, "")
	result := mb.String()

	assert.Contains(t, result, fromID, "Result should contain from node ID")
	assert.Contains(t, result, toID, "Result should contain to node ID")
}

func TestHashID(t *testing.T) {
	id1 := hashID("test")
	id2 := hashID("test")
	id3 := hashID("different")

	// Same input should produce same hash
	assert.Equal(t, id1, id2, "Same input should produce same hash")

	// Different input should produce different hash
	assert.NotEqual(t, id1, id3, "Different input should produce different hash")

	// Hash should be a valid hex string
	assert.Regexp(t, `^[0-9a-f]+$`, id1, "Hash should be hex string")
}

func TestEscapeMermaidLabel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "backslash",
			input:    "test\\path",
			expected: "test\\\\path",
		},
		{
			name:     "quotes",
			input:    `test"quote"`,
			expected: `test\"quote\"`,
		},
		{
			name:     "newline",
			input:    "test\nline",
			expected: "test line",
		},
		{
			name:     "carriage return",
			input:    "test\rline",
			expected: "test line",
		},
		{
			name:     "tab",
			input:    "test\tline",
			expected: "test line",
		},
		{
			name:     "multiple special chars",
			input:    "test\"\\\n\r\t",
			expected: "test\\\"\\\\   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeMermaidLabel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildMermaidFromTraces(t *testing.T) {
	now := time.Now()
	traces := []tracer.TraceEvent{
		{
			TS:       now,
			Local:    "127.0.0.1:8080",
			Remote:   "127.0.0.1:12345",
			Dir:      tracer.TraceIn,
			Len:      10,
			ClientID: 1,
		},
		{
			TS:       now.Add(time.Second),
			Local:    "127.0.0.1:8080",
			Remote:   "127.0.0.1:12345",
			Dir:      tracer.TraceOut,
			Len:      20,
			ClientID: 1,
		},
	}

	result := BuildMermaidFromTraces(traces)
	require.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, "flowchart TD"), "Should start with flowchart TD")
	assert.Contains(t, result, "local: 127.0.0.1:8080", "Should contain local address")
	assert.Contains(t, result, "remote: 127.0.0.1:12345", "Should contain remote address")
}

func TestBuildSequenceDiagramFromTraces(t *testing.T) {
	now := time.Now()
	traces := []tracer.TraceEvent{
		{
			TS:       now,
			Local:    "127.0.0.1:8080",
			Remote:   "127.0.0.1:12345",
			Dir:      tracer.TraceIn,
			Len:      10,
			ClientID: 1,
		},
		{
			TS:       now.Add(time.Second),
			Local:    "127.0.0.1:8080",
			Remote:   "127.0.0.1:12345",
			Dir:      tracer.TraceOut,
			Len:      20,
			ClientID: 1,
		},
	}

	result := BuildSequenceDiagramFromTraces(traces)
	require.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, "sequenceDiagram"), "Should start with sequenceDiagram")
	assert.Contains(t, result, "Server", "Should contain Server participant")
	assert.Contains(t, result, "Client", "Should contain Client participant")
}

func TestBuildSequenceDiagramFromTraces_Empty(t *testing.T) {
	result := BuildSequenceDiagramFromTraces([]tracer.TraceEvent{})
	require.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, "sequenceDiagram"), "Should start with sequenceDiagram even with empty traces")
}

func TestBuildMermaidFromTraces_Empty(t *testing.T) {
	result := BuildMermaidFromTraces([]tracer.TraceEvent{})
	require.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, "flowchart TD"), "Should start with flowchart TD even with empty traces")
}
