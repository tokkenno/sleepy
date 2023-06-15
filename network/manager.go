package network

import (
	"fmt"
	"net"
	"sleepy/network/common/udp"
)

type Manager interface {
	SendUDP(ip net.IP, port uint16, packet udp.Packet)
}

type manager struct {
}

var _ Manager = &manager{}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) SendUDP(ip net.IP, port uint16, packet udp.Packet) {
	fmt.Sprintf("Mandando paquete a %s:%d => ", ip.String(), port)
	fmt.Println(packet.GetData())
}
