package ws

import "golang.org/x/net/websocket"

func (wh *WebSocketHub) HandleWS(ws *websocket.Conn) {
	wh.mu.Lock()
	wh.conns[ws] = true
	wh.mu.Unlock()

	wh.readLoop(ws)

	wh.mu.Lock()
	delete(wh.conns, ws)
	wh.mu.Unlock()
}
