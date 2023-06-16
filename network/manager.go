package network

import (
	"fmt"
	"net"
	"sleepy/network/common/udp"
)

type Manager interface {
	SendUDP(ip net.IP, port uint16, packet udp.Packet) error
}

type manager struct {
}

var _ Manager = &manager{}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) SendUDP(ip net.IP, port uint16, packet udp.Packet) error {
	fmt.Sprintf("Mandando paquete a %s:%d => ", ip.String(), port)
	fmt.Println(packet.GetData())

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		return err
	}

	_, err = conn.Write(packet.GetData())
	if err != nil {
		return err
	}

	return nil
}
