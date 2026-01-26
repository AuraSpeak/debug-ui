package ws

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/websocket"
)

type WebsocketMessageType string

const (
	TypeLog WebsocketMessageType = "LOG"
)

type WebSocketMessage struct {
	Type    WebsocketMessageType `json:"type"`
	Content string               `json:"content"`
}

type WebSocketHub struct {
	conns  map[*websocket.Conn]bool
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewHub(ctx context.Context) *WebSocketHub {
	hubCtx, cancel := context.WithCancel(ctx)
	wsHub := &WebSocketHub{
		conns:  make(map[*websocket.Conn]bool),
		mu:     sync.Mutex{},
		ctx:    hubCtx,
		cancel: cancel,
	}
	log.AddHook(NewWebSocketHook(wsHub))
	return wsHub
}

func (wh *WebSocketHub) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-wh.ctx.Done():
			return
		default:
		}
		// Set the ReadDeadline, so we can peridoicly test the context
		err := ws.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		if err != nil {
			log.WithField("caller", "web").WithError(err).Error("readLoop error")
		}
		n, err := ws.Read(buf)
		if err != nil {
			// Client has Closed the connection
			if err == io.EOF {
				break
			}
			// Timeout-Fehler will be getting ignored, checked by the context
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.WithField("caller", "web").WithError(err).Error("readLoop error")
			break
		}
		msg := buf[:n]
		wh.Broadcast(msg)

	}
}

func (wh *WebSocketHub) Broadcast(b []byte) {
	// NOTE: Do NOT call log.Infof here! It would cause infinite recursion
	// because the WebSocketHook calls Broadcast, which would call log.Infof again
	wh.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(wh.conns))
	for ws := range wh.conns {
		conns = append(conns, ws)
	}
	wh.mu.Unlock()

	for _, ws := range conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				log.WithField("caller", "web").WithError(err).Error("broadcast error")
			}
		}(ws)
	}
}

// Cancel cancels the WebSocketHub context, signaling all goroutines to stop
func (wh *WebSocketHub) Cancel() {
	wh.cancel()
}

// Wait waits for all goroutines in the hub to finish
func (wh *WebSocketHub) Wait() {
	wh.wg.Wait()
}
