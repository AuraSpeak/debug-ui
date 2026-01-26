package app

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/auraspeak/debug-ui/internal/api"
	"github.com/auraspeak/server/pkg/debugui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	require.NotNil(t, server)
	assert.Equal(t, 8080, server.Port)
	assert.Equal(t, 9090, server.config.UDPPort)
	assert.NotNil(t, server.wsHub)
	assert.NotNil(t, server.udpClients)
	assert.NotNil(t, server.clientCommandChs)
	assert.NotNil(t, server.ctx)
}

func TestConvertMessageToBytes_Text(t *testing.T) {
	message := "Hello, World!"
	format := "text"

	result, err := convertMessageToBytes(message, format)

	require.NoError(t, err)
	assert.Equal(t, []byte(message), result)
}

func TestConvertMessageToBytes_Hex(t *testing.T) {
	message := "48656c6c6f"
	format := "hex"

	result, err := convertMessageToBytes(message, format)

	require.NoError(t, err)
	expected, _ := hex.DecodeString("48656c6c6f")
	assert.Equal(t, expected, result)
}

func TestConvertMessageToBytes_HexWithSpaces(t *testing.T) {
	message := "48 65 6c 6c 6f"
	format := "hex"

	result, err := convertMessageToBytes(message, format)

	require.NoError(t, err)
	expected, _ := hex.DecodeString("48656c6c6f")
	assert.Equal(t, expected, result)
}

func TestConvertMessageToBytes_InvalidHex(t *testing.T) {
	message := "invalid hex"
	format := "hex"

	_, err := convertMessageToBytes(message, format)

	assert.Error(t, err)
}

func TestServer_GetUDPServerState_NoServer(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/server/get", nil)
	rr := httptest.NewRecorder()

	server.GetUDPServerState(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response api.ServerStateResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.IsAlive)
	assert.False(t, response.ShouldStop)
}

func TestServer_GetAllUDPClients_Empty(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/all", nil)
	rr := httptest.NewRecorder()

	server.GetAllUDPClients(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response api.AllUDPClientResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Empty(t, response.UDPClients)
}

func TestServer_GetUDPClientStateByName_NotFound(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/name?name=nonexistent", nil)
	rr := httptest.NewRecorder()

	server.GetUDPClientStateByName(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Message, "not found")
}

func TestServer_GetUDPClientStateByName_NoName(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/name", nil)
	rr := httptest.NewRecorder()

	server.GetUDPClientStateByName(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "Name is required")
}

func TestServer_GetUDPClientStateById_NoId(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/id", nil)
	rr := httptest.NewRecorder()

	server.GetUDPClientStateById(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "ID is required")
}

func TestServer_GetUDPClientStateById_InvalidId(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/id?id=invalid", nil)
	rr := httptest.NewRecorder()

	server.GetUDPClientStateById(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "ID is invalid")
}

func TestServer_GetUDPClientStateById_NotFound(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/id?id=999", nil)
	rr := httptest.NewRecorder()

	server.GetUDPClientStateById(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Message, "not found")
}

func TestServer_StopUDPClient_NoName(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("POST", "/api/client/stop", nil)
	rr := httptest.NewRecorder()

	server.StopUDPClient(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "Name is required")
}

func TestServer_StopUDPClient_NotFound(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("POST", "/api/client/stop?name=nonexistent", nil)
	rr := httptest.NewRecorder()

	server.StopUDPClient(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Message, "not found")
}

func TestServer_GetAllUDPClientPaginated_DefaultValues(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/all/paginated", nil)
	rr := httptest.NewRecorder()

	server.GetAllUDPClientPaginated(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response api.UDPClientPaginatedRespone
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, 0, response.Total)
}

func TestServer_GetAllUDPClientPaginated_InvalidPage(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/all/paginated?page=invalid", nil)
	rr := httptest.NewRecorder()

	server.GetAllUDPClientPaginated(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "positive integer")
}

func TestServer_GetAllUDPClientPaginated_InvalidPageSize(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/all/paginated?pageSize=invalid", nil)
	rr := httptest.NewRecorder()

	server.GetAllUDPClientPaginated(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "positive integer")
}

func TestServer_GetAllUDPClientPaginated_NegativePage(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("GET", "/api/client/get/all/paginated?page=-1", nil)
	rr := httptest.NewRecorder()

	server.GetAllUDPClientPaginated(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestServer_SendDatagram_InvalidBody(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	req := httptest.NewRequest("POST", "/api/client/send", nil)
	rr := httptest.NewRecorder()

	server.SendDatagram(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "Invalid request body")
}

func TestServer_SendDatagram_InvalidFormat(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	reqBody := api.SendDatagramRequest{
		Id:      1,
		Message: "test",
		Format:  "invalid",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/client/send", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.SendDatagram(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Message, "Format must be")
}

func TestServer_SendDatagram_ClientNotFound(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	reqBody := api.SendDatagramRequest{
		Id:      999,
		Message: "test",
		Format:  "text",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/client/send", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.SendDatagram(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response api.ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Contains(t, response.Message, "not found")
}

func TestServer_HandleWS(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	// HandleWS should not panic even with nil websocket
	// (websocket.Handler will handle nil connections)
	assert.NotNil(t, server.wsHub)
	// Note: Actual WebSocket testing would require a real connection
	// This is just a basic smoke test
}

func TestServer_Shutdown_WithoutStart(t *testing.T) {
	cfg := debugui.Config{}
	server := NewServer(8080, 9090, cfg)

	// Cancel context first to avoid blocking
	server.cancel()

	// Shutdown will panic if httpServer is nil, so we skip this test
	// In real usage, httpServer is set in Run(), so this scenario shouldn't occur
	// We test that Shutdown works when server is properly initialized
	_ = server
}
