package api

import (
	"github.com/auraspeak/client"
)

type ActionType int

const (
	UDPClientActionAddDatagram ActionType = iota
)

type UDPClient struct {
	// ID of the client
	ID int
	// The UDP client
	Client *client.Client
	// name random chosen
	Name string
	// Datagram by the user
	Datagrams []Datagram
	// is it running
	Running bool
}
type UDPClientAction struct {
	ID     int
	Action ActionType
}
