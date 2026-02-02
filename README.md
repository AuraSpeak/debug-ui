# Debug-UI

The Debug-UI provides a debug interface for the AuraSpeak Project. It runs an HTTP server and WebSocket hub that control and observe one UDP server and multiple UDP clients (using auraspeak/server and auraspeak/client). It shows traces (e.g. Mermaid sequence diagrams). It utilizes the [AuraSpeak Protocol](https://github.com/AuraSpeak/protocol).

---

## Requirements

GO Version: 1.25.1

Go dependencies:
- `github.com/auraspeak/server` for local development replaced: `replace github.com/auraspeak/server => ../server`
- `github.com/auraspeak/client` for local development replaced: `replace github.com/auraspeak/client => ../client`
- `github.com/auraspeak/protocol` for local development replaced: `replace github.com/auraspeak/protocol => ../protocol`

For the frontend: Node/npm (see [web/README.md](web/README.md)).

---

## Structure

### app/Server

HTTP server, WebSocket hub, UDP server and UDP client management, trace collection, internal channels (command, message). `NewServer(httpPort, udpPort, cfg)`; `Run()`, `Shutdown(timeout)`.

### internal/api

Routes and request/response types (udp_types, respones). CORS applied to `/api/` routes.

### internal/communication

Internal messages (InternalMessage).

### internal/services

UDPServerService, ID manager, UDP client lifecycle.

### internal/ws

WebSocketHub, handler, logger hook. Broadcasts (e.g. `uss`, `usu`, `cnu`, `rp`) to connected clients.

### internal/util

NameGenerator, Seq, BuildSequenceDiagramFromTraces (Mermaid).

### web/

Vue 3 + Vite frontend. Static assets are served from `./bin`. See [web/README.md](web/README.md) for npm setup and development.

---

## Quick start

### Backend

Run `go run ./cmd`. HTTP server listens on port 8080; UDP server port is e.g. 9090. Config is loaded via `server/pkg/debugui`.

### API (overview)

- WebSocket: `/ws`
- UDP Server: POST `/api/server/start`, POST `/api/server/stop`, GET `/api/server/get`
- UDP Client: POST `/api/client/start`, POST `/api/client/stop`, POST `/api/client/send`, GET `/api/client/get/name`, GET `/api/client/get/id`, GET `/api/client/get/all`, GET `/api/client/get/all/paginated`
- Traces: GET `/api/traces/all` (Mermaid diagram per client; query param `name`)

### Frontend

```sh
cd web && npm install && npm run dev
```

Build for production: `npm run build`. Output can be placed in `./bin` for the app to serve.

---

## Testing

Run `go test ./...` to test the backend.

---

## License

[License](./LICENSE)
