package app

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/auraspeak/client"
	"github.com/auraspeak/client/pkg/command"
	"github.com/auraspeak/debug-ui/internal/api"
	"github.com/auraspeak/debug-ui/internal/communication"
	"github.com/auraspeak/debug-ui/internal/services"
	"github.com/auraspeak/debug-ui/internal/util"
	"github.com/auraspeak/debug-ui/internal/ws"
	"github.com/auraspeak/protocol"
	"github.com/auraspeak/server"
	serverCommand "github.com/auraspeak/server/pkg/command"
	serverConfig "github.com/auraspeak/server/pkg/debugui"
	"github.com/auraspeak/server/pkg/tracer"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Config struct {
	UDPPort int
}

type Server struct {
	Port       int
	mu         sync.Mutex
	config     Config
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
	shutdownWg sync.WaitGroup

	cfg *serverConfig.Config

	// WebSocket Hub
	wsHub *ws.WebSocketHub

	// UDP Parts
	udpServer *server.Server
	// udpClientWrapper
	udpClients      map[string]api.UDPClient
	clientMu        sync.Mutex
	udpClientAction api.UDPClientAction
	// Communicate from UDP Server and Clients to WebSocket Hub
	messageCh chan []communication.InternalMessage
	// Client command channels mapped by client ID
	clientCommandChs map[int]chan command.InternalCommand

	// Traces
	traces  []tracer.TraceEvent
	traceMu sync.Mutex
}

func NewServer(port int, udpPort int, cfg serverConfig.Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		Port:             port,
		mu:               sync.Mutex{},
		ctx:              ctx,
		cancel:           cancel,
		wsHub:            ws.NewHub(ctx),
		config:           Config{UDPPort: udpPort},
		udpClients:       make(map[string]api.UDPClient),
		clientCommandChs: make(map[int]chan command.InternalCommand),
		traceMu:          sync.Mutex{},
		cfg:              &cfg,
	}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", s.Port),
		Handler: api.RegisterRoutes(
			func(w http.ResponseWriter, r *http.Request) {
				// WebSocket handler needs special handling
				websocket.Handler(s.HandleWS).ServeHTTP(w, r)
			},
			s.StartUDPServer,
			s.StopUDPServer,
			s.GetUDPServerState,
			s.StartUDPClient,
			s.StopUDPClient,
			s.SendDatagram,
			s.GetUDPClientStateByName,
			s.GetUDPClientStateById,
			s.GetAllUDPClients,
			s.GetTraces,
			s.GetAllUDPClientPaginated,
		),
	}

	s.shutdownWg.Go(func() {
		s.handleInternal()
	})

	s.shutdownWg.Go(func() {
		s.handleTrace()
	})

	fmt.Printf("Starting server on http://localhost:%d\n", s.Port)
	// Broadcast restart signal once to all clients
	s.wsHub.Broadcast([]byte("rp"))
	return s.httpServer.ListenAndServe()
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	if s.wsHub != nil {
		s.wsHub.HandleWS(ws)
	}
}

// handleClientCommands listens for commands from a specific UDP client
func (s *Server) handleClientCommands(clientID int, cmdCh chan command.InternalCommand) {
	s.shutdownWg.Go(func() {
		for {
			select {
			case cmd := <-cmdCh:
				switch cmd {
				case command.CmdUpdateClientState:
					s.mu.Lock()
					// Find udpClient by ID and update running field
					for name, uc := range s.udpClients {
						if uc.ID == clientID {
							// Update running field from ClientState
							running := uc.Client.ClientState.Running
							uc.Running = running == 1
							// Update the map entry
							s.udpClients[name] = uc

							// Broadcast to all WebSocket Clients that the UDP Client State has changed
							if s.wsHub != nil {
								s.wsHub.Broadcast([]byte("usu" + strconv.Itoa(clientID)))
							}
							break
						}
					}
					s.mu.Unlock()
				}
			case <-s.ctx.Done():
				return
			}
		}
	})
}

// Handles all internal communications to the web server
// Avalible Commands:
// uss: tells the clients, that the server state has been updated
// usu: tells the clients, that a udp client state has been updated
// cnu: tells the web server, that a new udp client has been started
func (s *Server) handleInternal() {
	s.shutdownWg.Go(func() {
		// Brodcast through on UDP Server State Changes
		for {
			// Wait until UDP server is initialized
			s.mu.Lock()
			udpServer := s.udpServer
			s.mu.Unlock()

			if udpServer == nil {
				// UDP server not started yet, wait a bit and check again
				select {
				case <-s.ctx.Done():
					return
				case <-time.After(100 * time.Millisecond):
					continue
				}
			}

			// UDP server is available, listen for commands
			select {
			case cmd := <-udpServer.OutCommandCh:
				switch cmd {
				case serverCommand.CmdUpdateServerState:
					// Broadcast to all WebSocket Clients that the UDP Server State has changed
					if s.wsHub != nil {
						s.wsHub.Broadcast([]byte("uss"))
					}
				}
			case <-s.ctx.Done():
				return
			}
		}
	})
}

func (s *Server) handleTrace() {
	s.shutdownWg.Go(func() {
		for {
			s.mu.Lock()
			udpServer := s.udpServer
			s.mu.Unlock()

			if udpServer == nil {
				select {
				case <-s.ctx.Done():
					return
				case <-time.After(100 * time.Millisecond):
					continue
				}
			}

			select {
			case <-s.ctx.Done():
				return
			case trace := <-s.udpServer.TraceCh:
				s.traceMu.Lock()
				s.traces = append(s.traces, trace)
				log.WithFields(log.Fields{
					"caller": "web",
					"cid":    trace.ClientID,
				}).Debugf("Received trace: %+v", trace)
				s.traceMu.Unlock()
			}
		}
	})
}

func (s *Server) Shutdown(timeout time.Duration) error {
	fmt.Println("Shutting down server...")

	// Signals all go routines to cancel
	s.cancel()

	// Cancel the WebSocketHub
	if s.wsHub != nil {
		s.wsHub.Cancel()
	}

	// Create a context with the given timeout for all shutdown operations
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shutdown HTTP server
	httpDone := make(chan error, 1)
	go func() {
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			httpDone <- fmt.Errorf("error shutting down HTTP server: %w", err)
		} else {
			httpDone <- nil
		}
	}()

	// Wait for all goroutines to finish or timeout
	wgDone := make(chan struct{})
	go func() {
		s.shutdownWg.Wait()
		close(wgDone)
	}()

	// Wait for either all goroutines or context timeout
	select {
	case <-wgDone:
		fmt.Println("All connections closed")
	case <-shutdownCtx.Done():
		fmt.Println("Warning: Shutdown timeout reached, some connections may not have closed gracefully")
	}

	// Wait for HTTP server shutdown result with timeout
	select {
	case err := <-httpDone:
		if err != nil {
			return err
		}
	case <-shutdownCtx.Done():
		fmt.Println("HTTP server shutdown timeout")
	}

	fmt.Println("Server shutdown complete")
	return nil
}

// Helper functions for UDP client management

// genUDPClient creates a new UDP client and returns its name
func (s *Server) genUDPClient(port int) string {
	name := util.GetFirstName()
	id := services.GetNextID()
	client := client.NewDebugClient("localhost", port, id)
	s.mu.Lock()
	s.udpClients[name] = api.UDPClient{
		ID:        id,
		Client:    client,
		Name:      name,
		Datagrams: []api.Datagram{},
		Running:   false,
	}
	// Register client command channel and start listening
	s.clientCommandChs[id] = client.OutCommandCh
	s.handleClientCommands(id, client.OutCommandCh)
	s.mu.Unlock()
	log.Infof("UDP client started: %s with id %d", name, id)
	return name
}

// convertMessageToBytes converts a message string to []byte based on format
func convertMessageToBytes(message string, format string) ([]byte, error) {
	if format == "hex" {
		// Remove spaces from hex string
		hexString := strings.ReplaceAll(strings.TrimSpace(message), " ", "")
		return hex.DecodeString(hexString)
	}
	// Text format: convert string to []byte using UTF-8
	return []byte(message), nil
}

// handleAllClient handles all incoming packets from UDP clients
func (s *Server) handleAllClient(name string, packet []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	udpClient, ok := s.udpClients[name]
	if !ok {
		log.Errorf("UDP client not found: %s", name)
		return fmt.Errorf("UDP client not found: %s", name)
	}

	newDatagram := api.Datagram{
		Direction: api.ServerToClient,
		Message:   packet,
	}
	udpClient.Datagrams = append(udpClient.Datagrams, newDatagram)
	s.udpClients[name] = udpClient
	if s.wsHub != nil {
		s.wsHub.Broadcast([]byte("usu" + strconv.Itoa(udpClient.ID)))
	}

	return nil
}

// UDP Client Handler Methods

func (s *Server) StartUDPClient(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	name := s.genUDPClient(s.config.UDPPort)
	udpClient := s.udpClients[name]
	udpClient.Client.OnPacket(protocol.PacketTypeDebugAny, func(packet *protocol.Packet) error {
		return s.handleAllClient(name, packet.Payload)
	})
	s.mu.Unlock()
	s.shutdownWg.Go(func() {
		udpClient.Client.Run()
	})
	udpClientResponse := api.UDPClientResponse{
		Name: name,
		Id:   udpClient.ID,
	}

	s.wsHub.Broadcast([]byte("cnu"))
	udpClientResponse.Send(w)
}

func (s *Server) StopUDPClient(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	name := r.URL.Query().Get("name")
	if name == "" {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
		apiError.Send(w)
		return
	}
	udpClient, ok := s.udpClients[name]
	if !ok {
		apiError := api.ApiError{
			Code:    http.StatusNotFound,
			Message: "UDP client not found",
		}
		apiError.Send(w)
		return
	}
	udpClient.Client.Stop()
	// Remove client command channel from map
	delete(s.clientCommandChs, udpClient.ID)
	apiSuccess := api.ApiSuccess{
		Message: "UDP client stopped",
	}
	apiSuccess.Send(w)
}

func (s *Server) GetUDPClientStateByName(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	name := r.URL.Query().Get("name")
	if name == "" {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
		apiError.Send(w)
		return
	}
	udpClient, ok := s.udpClients[name]
	if !ok {
		apiError := api.ApiError{
			Code:    http.StatusNotFound,
			Message: "UDP client not found",
		}
		apiError.Send(w)
		return
	}
	udpClientStateResponse := api.UDPClientStateResponse{
		Id:        udpClient.ID,
		Running:   udpClient.Running,
		Datagrams: udpClient.Datagrams,
	}
	udpClientStateResponse.Send(w)
}

func (s *Server) GetUDPClientStateById(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := r.URL.Query().Get("id")
	if id == "" {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "ID is required",
		}
		apiError.Send(w)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "ID is invalid",
		}
		apiError.Send(w)
		return
	}
	for _, udpClient := range s.udpClients {
		if udpClient.ID == idInt {
			udpClientStateResponse := api.UDPClientStateResponse{
				Id:        udpClient.ID,
				Running:   udpClient.Running,
				Datagrams: udpClient.Datagrams,
			}
			udpClientStateResponse.Send(w)
			return
		}
	}
	apiError := api.ApiError{
		Code:    http.StatusNotFound,
		Message: "UDP client not found",
	}
	apiError.Send(w)
}

func (s *Server) GetAllUDPClients(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	allUDPClientResponse := api.AllUDPClientResponse{
		UDPClients: []api.UDPClientResponse{},
	}
	for name, udpClient := range s.udpClients {
		allUDPClientResponse.UDPClients = append(allUDPClientResponse.UDPClients, api.UDPClientResponse{
			Id:   udpClient.ID,
			Name: name,
		})
	}
	allUDPClientResponse.Send(w)
}

func (s *Server) GetAllUDPClientPaginated(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse page parameter with default value 1
	pageStr := r.URL.Query().Get("page")
	pageInt := 1
	if pageStr != "" {
		var err error
		pageInt, err = strconv.Atoi(pageStr)
		if err != nil || pageInt < 1 {
			apiError := api.ApiError{
				Code:    http.StatusBadRequest,
				Message: "Page must be a positive integer",
			}
			apiError.Send(w)
			return
		}
	}

	// Parse pageSize parameter with default value 10
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSizeInt := 10
	if pageSizeStr != "" {
		var err error
		pageSizeInt, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSizeInt < 1 {
			apiError := api.ApiError{
				Code:    http.StatusBadRequest,
				Message: "Page size must be a positive integer",
			}
			apiError.Send(w)
			return
		}
	}

	// Parse search query parameter
	searchQuery := r.URL.Query().Get("q")

	// Filter and collect all matching UDP clients (only basic info: id and name)
	allItems := []api.UDPClientListItem{}
	for name, udpClient := range s.udpClients {
		// Apply search filter if provided
		if searchQuery != "" {
			// Case-insensitive search in client name
			if !strings.Contains(strings.ToLower(name), strings.ToLower(searchQuery)) {
				continue
			}
		}

		allItems = append(allItems, api.UDPClientListItem{
			Id:   udpClient.ID,
			Name: name,
		})
	}

	// Calculate pagination
	total := len(allItems)
	totalPages := int(math.Ceil(float64(total) / float64(pageSizeInt)))

	// Validate page number
	if pageInt > totalPages && totalPages > 0 {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Page is out of range",
		}
		apiError.Send(w)
		return
	}

	// Calculate slice indices for pagination
	startIndex := (pageInt - 1) * pageSizeInt
	endIndex := startIndex + pageSizeInt
	if endIndex > total {
		endIndex = total
	}

	// Extract items for current page
	var items []api.UDPClientListItem
	if startIndex < total {
		items = allItems[startIndex:endIndex]
	} else {
		items = []api.UDPClientListItem{}
	}

	// Create and send response
	paginatedResponse := api.UDPClientPaginatedRespone{
		Items:    items,
		Page:     pageInt,
		PageSize: pageSizeInt,
		Total:    total,
	}
	paginatedResponse.Send(w)
}

func (s *Server) SendDatagram(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req api.SendDatagramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		}
		apiError.Send(w)
		return
	}

	// Validate format
	if req.Format != "hex" && req.Format != "text" {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Format must be 'hex' or 'text'",
		}
		apiError.Send(w)
		return
	}

	// Find client by ID and validate (with lock)
	s.mu.Lock()
	var clientToSend *client.Client
	var clientName string
	var clientRunning bool
	for name, uc := range s.udpClients {
		if uc.ID == req.Id {
			clientToSend = uc.Client
			clientName = name
			clientRunning = uc.Running
			break
		}
	}
	s.mu.Unlock()

	if clientToSend == nil {
		apiError := api.ApiError{
			Code:    http.StatusNotFound,
			Message: "UDP client not found",
		}
		apiError.Send(w)
		return
	}

	// Check if client is running
	if !clientRunning {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Client is not running",
		}
		apiError.Send(w)
		return
	}

	// Convert message to []byte based on format
	messageBytes, err := convertMessageToBytes(req.Message, req.Format)
	if err != nil {
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "Invalid hex string",
			Details: err.Error(),
		}
		apiError.Send(w)
		return
	}

	packet := &protocol.Packet{
		PacketHeader: protocol.Header{PacketType: protocol.PacketTypeDebugAny},
		Payload:      messageBytes,
	}

	// Send message via client (außerhalb des Locks, damit es nicht blockiert)
	if err := clientToSend.Send(packet.Encode()); err != nil {
		apiError := api.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to send datagram",
			Details: err.Error(),
		}
		apiError.Send(w)
		return
	}

	// Lock wieder holen für Map-Update
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find client again (könnte sich geändert haben)
	udpClient, ok := s.udpClients[clientName]
	if !ok {
		apiError := api.ApiError{
			Code:    http.StatusNotFound,
			Message: "UDP client not found",
		}
		apiError.Send(w)
		return
	}

	// Store datagram in client's datagrams list
	newDatagram := api.Datagram{
		Direction: api.ClientToServer,
		Message:   messageBytes,
	}
	udpClient.Datagrams = append(udpClient.Datagrams, newDatagram)
	s.udpClients[clientName] = udpClient

	// Broadcast WebSocket update
	if s.wsHub != nil {
		s.wsHub.Broadcast([]byte("usu" + strconv.Itoa(req.Id)))
	}

	// Send success response
	response := api.SendDatagramResponse{
		Message: "Datagram sent successfully",
	}
	log.Infof("Datagram sent successfully: %s", string(messageBytes))
	response.Send(w)
}

func (s *Server) GetTraces(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	clientName := r.URL.Query().Get("name")
	udpClient, ok := s.udpClients[clientName]
	if !ok {
		apiError := api.ApiError{
			Code:    http.StatusNotFound,
			Message: "UDP client not found",
		}
		apiError.Send(w)
		return
	}
	clientID := udpClient.ID
	s.traceMu.Lock()
	traces := s.traces
	s.traceMu.Unlock()
	filteredTraces := []tracer.TraceEvent{}
	for _, trace := range traces {
		if trace.ClientID == clientID {
			filteredTraces = append(filteredTraces, trace)
		}
	}
	md := util.BuildSequenceDiagramFromTraces(filteredTraces)
	traceRes := api.MermaidResponse{
		Heading: fmt.Sprintf("Diagram for user: %s", clientName),
		Diagram: md,
	}
	traceRes.Send(w)
}

// UDP Server Handler Methods

func (s *Server) StartUDPServer(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	udpServerService := services.NewUDPServerService(s.ctx, s.cfg, s.wsHub)
	if err := udpServerService.Start(s.config.UDPPort); err != nil {
		apiError := api.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create server (DTLS config error)",
		}
		apiError.Send(w)
		return
	}

	udpServer := udpServerService.GetServer()
	if udpServer == nil {
		apiError := api.ApiError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create server (DTLS config error)",
		}
		apiError.Send(w)
		return
	}

	s.udpServer = udpServer
	udpServer.OnPacket(protocol.PacketTypeDebugAny, func(packet *protocol.Packet, clientAddr string) error {
		return udpServerService.HandleAll(clientAddr, packet.Payload)
	})

	s.shutdownWg.Go(func() {
		if err := udpServer.Run(); err != nil {
			log.WithField("caller", "web").WithError(err).Error("error starting udp server")
		}
	})

	apiSuccess := api.ApiSuccess{
		Message: "UDP server started",
	}
	apiSuccess.Send(w)
}

func (s *Server) StopUDPServer(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	udpServer := s.udpServer
	s.mu.Unlock()

	if udpServer == nil {
		log.WithField("caller", "web").Warn("UDP server is not running")
		apiError := api.ApiError{
			Code:    http.StatusBadRequest,
			Message: "UDP server is not running",
		}
		apiError.Send(w)
		return
	}

	udpServer.Stop()
	log.WithField("caller", "web").Info("UDP server stopped")
	apiSuccess := api.ApiSuccess{
		Message: "UDP server stopped",
	}
	apiSuccess.Send(w)
}

func (s *Server) GetUDPServerState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	udpServer := s.udpServer
	s.mu.Unlock()

	if udpServer == nil {
		serverStateResponse := api.ServerStateResponse{
			ShouldStop: false,
			IsAlive:    false,
		}
		serverStateResponse.Send(w)
		return
	}

	state := udpServer.ServerState
	serverStateResponse := api.ServerStateResponse{
		ShouldStop: state.ShouldStop,
		IsAlive:    state.IsAlive,
	}
	serverStateResponse.Send(w)
}
