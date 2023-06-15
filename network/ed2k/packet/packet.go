package packet

import (
	"sleepy/network/common/udp"
	"sleepy/network/ed2k/common"
)

type Packet struct {
	udp.RawPacket
}

func NewPacket(
	protocol common.Protocol,
	size int,
) *Packet {
	packet := &Packet{}
	if protocol == common.ProtocolEd2kServerUDP {
		packet.RawPacket = *udp.NewPacket(size + 2)
	} else if protocol == common.ProtocolEd2kPeerUDP {
		packet.RawPacket = *udp.NewPacket(size + 4)
	} else {
		packet.RawPacket = *udp.NewPacket(size)
	}
	packet.SetProtocol(common.ProtocolEd2kServerUDP)
	return packet
}

func (packet *Packet) SetProtocol(protocol common.Protocol) {
	packet.SetByte(0, byte(protocol))
}

func (packet *Packet) GetProtocol() common.Protocol {
	return common.Protocol(packet.GetByte(0))
}

func (packet *Packet) SetCommand(command byte) {
	if packet.GetProtocol() == common.ProtocolEd2kServerUDP {
		packet.SetByte(1, command)
	} else {
		packet.SetByte(5, command)
	}
}

func (packet *Packet) GetCommand() byte {
	if packet.GetProtocol() == common.ProtocolEd2kServerUDP {
		return packet.GetByte(1)
	} else {
		return packet.GetByte(5)
	}
}
