package services

import (
	"context"
	"errors"
	"sync"

	"github.com/auraspeak/debug-ui/internal/ws"
	"github.com/auraspeak/protocol"
	"github.com/auraspeak/server"
	"github.com/auraspeak/server/pkg/debugui"
	log "github.com/sirupsen/logrus"
)

type UDPServerService struct {
	server *server.Server
	mu     sync.Mutex
	ctx    context.Context
	cfg    *debugui.Config
	wsHub  *ws.WebSocketHub
}

func NewUDPServerService(ctx context.Context, cfg *debugui.Config, wsHub *ws.WebSocketHub) *UDPServerService {
	return &UDPServerService{
		server: nil,
		mu:     sync.Mutex{},
		ctx:    ctx,
		cfg:    cfg,
		wsHub:  wsHub,
	}
}

func (s *UDPServerService) Start(port int) error {
	s.server = server.NewServer(port, s.ctx, s.cfg)
	if s.server == nil {
		return errors.New("failed to create server")
	}
	return nil
}

func (s *UDPServerService) GetServer() *server.Server {
	return s.server
}

// handleAll handles all incoming packets from UDP server
func (s *UDPServerService) HandleAll(clientAddr string, packet []byte) error {
	log.WithField("caller", "web").Infof("Received packet: %s", string(packet))
	s.mu.Lock()
	if s.server != nil {
		s.server.Broadcast(&protocol.Packet{
			PacketHeader: protocol.Header{PacketType: protocol.PacketTypeDebugAny},
			Payload:      packet,
		})
	}
	s.mu.Unlock()
	s.mu.Lock()
	if s.wsHub != nil {
		s.wsHub.Broadcast(packet)
	}
	s.mu.Unlock()
	return nil
}
