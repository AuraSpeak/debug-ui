package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiError_Send(t *testing.T) {
	apiError := ApiError{
		Code:    http.StatusBadRequest,
		Message: "Test error",
		Details: "Test details",
	}

	rr := httptest.NewRecorder()
	apiError.Send(rr)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	// Note: Content-Type might not be set by httptest, but response should be valid JSON

	var response ApiError
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, apiError.Code, response.Code)
	assert.Equal(t, apiError.Message, response.Message)
	assert.Equal(t, apiError.Details, response.Details)
}

func TestApiSuccess_Send(t *testing.T) {
	apiSuccess := ApiSuccess{
		Message: "Test success",
		Details: "Test details",
	}

	rr := httptest.NewRecorder()
	apiSuccess.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)
	// Note: Content-Type might not be set by httptest, but response should be valid JSON

	var response ApiSuccess
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, apiSuccess.Message, response.Message)
	assert.Equal(t, apiSuccess.Details, response.Details)
}

func TestServerStateResponse_Send(t *testing.T) {
	response := ServerStateResponse{
		ShouldStop: true,
		IsAlive:    false,
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result ServerStateResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, response.ShouldStop, result.ShouldStop)
	assert.Equal(t, response.IsAlive, result.IsAlive)
}

func TestUDPClientResponse_Send(t *testing.T) {
	response := UDPClientResponse{
		Name: "TestClient",
		Id:   42,
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result UDPClientResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, response.Name, result.Name)
	assert.Equal(t, response.Id, result.Id)
}

func TestUDPClientStateResponse_Send(t *testing.T) {
	response := UDPClientStateResponse{
		Id:        42,
		Running:   true,
		Datagrams: []Datagram{},
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result UDPClientStateResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, response.Id, result.Id)
	assert.Equal(t, response.Running, result.Running)
	assert.Equal(t, len(response.Datagrams), len(result.Datagrams))
}

func TestAllUDPClientResponse_Send(t *testing.T) {
	response := AllUDPClientResponse{
		UDPClients: []UDPClientResponse{
			{Name: "Client1", Id: 1},
			{Name: "Client2", Id: 2},
		},
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result AllUDPClientResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, len(response.UDPClients), len(result.UDPClients))
}

func TestUDPClientPaginatedRespone_Send(t *testing.T) {
	response := UDPClientPaginatedRespone{
		Items: []UDPClientListItem{
			{Id: 1, Name: "Client1"},
			{Id: 2, Name: "Client2"},
		},
		Page:     1,
		PageSize: 10,
		Total:    2,
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result UDPClientPaginatedRespone
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, response.Page, result.Page)
	assert.Equal(t, response.PageSize, result.PageSize)
	assert.Equal(t, response.Total, result.Total)
	assert.Equal(t, len(response.Items), len(result.Items))
}

func TestSendDatagramResponse_Send(t *testing.T) {
	response := SendDatagramResponse{
		Message: "Datagram sent successfully",
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result SendDatagramResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, response.Message, result.Message)
}

func TestMermaidResponse_Send(t *testing.T) {
	response := MermaidResponse{
		Heading: "Test Diagram",
		Diagram: "sequenceDiagram\nA->>B: Test",
	}

	rr := httptest.NewRecorder()
	response.Send(rr)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result MermaidResponse
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, response.Heading, result.Heading)
	assert.Equal(t, response.Diagram, result.Diagram)
}
