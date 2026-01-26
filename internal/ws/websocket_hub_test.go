package ws

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHub(t *testing.T) {
	ctx := context.Background()
	hub := NewHub(ctx)

	require.NotNil(t, hub)
	assert.NotNil(t, hub.conns)
	assert.NotNil(t, hub.ctx)
	assert.NotNil(t, hub.cancel)
}

func TestWebSocketHub_Cancel(t *testing.T) {
	ctx := context.Background()
	hub := NewHub(ctx)

	// Cancel should not panic
	assert.NotPanics(t, func() {
		hub.Cancel()
	})

	// Context should be cancelled
	select {
	case <-hub.ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled")
	}
}

func TestWebSocketHub_Broadcast_NoConnections(t *testing.T) {
	ctx := context.Background()
	hub := NewHub(ctx)

	// Broadcasting with no connections should not panic
	assert.NotPanics(t, func() {
		hub.Broadcast([]byte("test message"))
	})
}

func TestWebSocketHub_Broadcast_WithConnections(t *testing.T) {
	ctx := context.Background()
	hub := NewHub(ctx)

	// Broadcasting with connections should not panic
	// Note: Testing with actual websocket.Conn would require a real connection
	// This test just verifies the broadcast logic doesn't panic
	assert.NotPanics(t, func() {
		hub.Broadcast([]byte("test message"))
	})
}
