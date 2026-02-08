package api

import (
	"net/http"

	"github.com/auraspeak/debug-ui/internal/middleware"
)

// RegisterRoutes creates an HTTP handler with all API routes
// Handler functions are passed as parameters to avoid import cycles
func RegisterRoutes(
	handleWS http.HandlerFunc,
	startUDPServer http.HandlerFunc,
	stopUDPServer http.HandlerFunc,
	getUDPServerState http.HandlerFunc,
	startUDPClient http.HandlerFunc,
	stopUDPClient http.HandlerFunc,
	sendDatagram http.HandlerFunc,
	getUDPClientStateByName http.HandlerFunc,
	getUDPClientStateById http.HandlerFunc,
	getAllUDPClients http.HandlerFunc,
	getTraces http.HandlerFunc,
	getAllUDPClientPaginated http.HandlerFunc,
	getClientMap http.HandlerFunc,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWS)
	mux.Handle("/", http.FileServer(http.Dir("./bin")))

	// UDP Server handlers
	mux.HandleFunc("POST /api/server/start", startUDPServer)
	mux.HandleFunc("POST /api/server/stop", stopUDPServer)
	mux.HandleFunc("GET /api/server/get", getUDPServerState)

	mux.HandleFunc("POST /api/client/start", startUDPClient)
	mux.HandleFunc("POST /api/client/stop", stopUDPClient)
	mux.HandleFunc("POST /api/client/send", sendDatagram)
	mux.HandleFunc("GET /api/client/get/name", getUDPClientStateByName)
	mux.HandleFunc("GET /api/client/get/id", getUDPClientStateById)
	mux.HandleFunc("GET /api/client/get/all", getAllUDPClients)
	mux.HandleFunc("GET /api/client/map", getClientMap)

	// Trace handlers
	mux.HandleFunc("GET /api/traces/all", getTraces)
	// Paginated all UDP clients
	mux.HandleFunc("GET /api/client/get/all/paginated", getAllUDPClientPaginated)

	// Wrap the entire mux with CORS for all /api/ routes
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply CORS to /api/ routes
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			middleware.CorsWrapper(mux).ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})
}
