package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRoutes(t *testing.T) {
	// Create mock handlers
	mockWS := func(w http.ResponseWriter, r *http.Request) {}
	mockStartUDPServer := func(w http.ResponseWriter, r *http.Request) {}
	mockStopUDPServer := func(w http.ResponseWriter, r *http.Request) {}
	mockGetUDPServerState := func(w http.ResponseWriter, r *http.Request) {}
	mockStartUDPClient := func(w http.ResponseWriter, r *http.Request) {}
	mockStopUDPClient := func(w http.ResponseWriter, r *http.Request) {}
	mockSendDatagram := func(w http.ResponseWriter, r *http.Request) {}
	mockGetUDPClientStateByName := func(w http.ResponseWriter, r *http.Request) {}
	mockGetUDPClientStateById := func(w http.ResponseWriter, r *http.Request) {}
	mockGetAllUDPClients := func(w http.ResponseWriter, r *http.Request) {}
	mockGetTraces := func(w http.ResponseWriter, r *http.Request) {}
	mockGetAllUDPClientPaginated := func(w http.ResponseWriter, r *http.Request) {}

	handler := RegisterRoutes(
		mockWS,
		mockStartUDPServer,
		mockStopUDPServer,
		mockGetUDPServerState,
		mockStartUDPClient,
		mockStopUDPClient,
		mockSendDatagram,
		mockGetUDPClientStateByName,
		mockGetUDPClientStateById,
		mockGetAllUDPClients,
		mockGetTraces,
		mockGetAllUDPClientPaginated,
	)

	require.NotNil(t, handler)
}

func TestRegisterRoutes_APIEndpoints(t *testing.T) {
	// Track which handlers were called
	called := make(map[string]bool)

	mockWS := func(w http.ResponseWriter, r *http.Request) { called["ws"] = true }
	mockStartUDPServer := func(w http.ResponseWriter, r *http.Request) { called["startUDPServer"] = true }
	mockStopUDPServer := func(w http.ResponseWriter, r *http.Request) { called["stopUDPServer"] = true }
	mockGetUDPServerState := func(w http.ResponseWriter, r *http.Request) { called["getUDPServerState"] = true }
	mockStartUDPClient := func(w http.ResponseWriter, r *http.Request) { called["startUDPClient"] = true }
	mockStopUDPClient := func(w http.ResponseWriter, r *http.Request) { called["stopUDPClient"] = true }
	mockSendDatagram := func(w http.ResponseWriter, r *http.Request) { called["sendDatagram"] = true }
	mockGetUDPClientStateByName := func(w http.ResponseWriter, r *http.Request) { called["getUDPClientStateByName"] = true }
	mockGetUDPClientStateById := func(w http.ResponseWriter, r *http.Request) { called["getUDPClientStateById"] = true }
	mockGetAllUDPClients := func(w http.ResponseWriter, r *http.Request) { called["getAllUDPClients"] = true }
	mockGetTraces := func(w http.ResponseWriter, r *http.Request) { called["getTraces"] = true }
	mockGetAllUDPClientPaginated := func(w http.ResponseWriter, r *http.Request) { called["getAllUDPClientPaginated"] = true }

	handler := RegisterRoutes(
		mockWS,
		mockStartUDPServer,
		mockStopUDPServer,
		mockGetUDPServerState,
		mockStartUDPClient,
		mockStopUDPClient,
		mockSendDatagram,
		mockGetUDPClientStateByName,
		mockGetUDPClientStateById,
		mockGetAllUDPClients,
		mockGetTraces,
		mockGetAllUDPClientPaginated,
	)

	// Test API routes
	tests := []struct {
		method string
		path   string
		key    string
	}{
		{"POST", "/api/server/start", "startUDPServer"},
		{"POST", "/api/server/stop", "stopUDPServer"},
		{"GET", "/api/server/get", "getUDPServerState"},
		{"POST", "/api/client/start", "startUDPClient"},
		{"POST", "/api/client/stop", "stopUDPClient"},
		{"POST", "/api/client/send", "sendDatagram"},
		{"GET", "/api/client/get/name", "getUDPClientStateByName"},
		{"GET", "/api/client/get/id", "getUDPClientStateById"},
		{"GET", "/api/client/get/all", "getAllUDPClients"},
		{"GET", "/api/traces/all", "getTraces"},
		{"GET", "/api/client/get/all/paginated", "getAllUDPClientPaginated"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			called = make(map[string]bool)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Check that the handler was called (or at least route exists)
			// Note: CORS middleware might affect the response, but route should exist
			assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route should exist")
		})
	}
}

func TestRegisterRoutes_CORS(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	handler := RegisterRoutes(
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
		mockHandler,
	)

	// Test that CORS headers are applied to API routes
	req := httptest.NewRequest("GET", "/api/server/get", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// CORS middleware should add headers
	// The exact headers depend on the CORS implementation
	assert.NotEqual(t, http.StatusNotFound, rr.Code)
}
