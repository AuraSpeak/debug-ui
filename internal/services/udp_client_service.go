package services

import (
	"context"
	"sync"

	"github.com/auraspeak/client"
	"github.com/auraspeak/debug-ui/internal/ws"
	"github.com/auraspeak/server/pkg/debugui"
)

type UDPClientService struct {
	mu               sync.Mutex
	udpClientService *client.Client
	ctx              context.Context
	cfg              *debugui.Config
	wsHub            *ws.WebSocketHub
	udpPort          int
	wg               sync.WaitGroup
}

func NewUDPClientService(ctx context.Context, cfg *debugui.Config, wsHub *ws.WebSocketHub, udpPort int) *UDPClientService {
	return &UDPClientService{
		mu:               sync.Mutex{},
		udpClientService: nil,
		ctx:              ctx,
		cfg:              cfg,
		wsHub:            wsHub,
		udpPort:          udpPort,
		wg:               sync.WaitGroup{},
	}
}

func NewUDPClient(ctx context.Context, cfg *debugui.Config, wsHub *ws.WebSocketHub, udpPort int) *UDPClientService {
	return &UDPClientService{
		mu:               sync.Mutex{},
		udpClientService: client.NewDebugClient("localhost", udpPort, 0),
		ctx:              ctx,
		cfg:              cfg,
		wsHub:            wsHub,
	}
}
