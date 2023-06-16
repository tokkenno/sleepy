package packet

import (
	"fmt"
	ed2kCommon "sleepy/network/ed2k/common"
	ed2kPacket "sleepy/network/ed2k/packet"
	"sleepy/network/kad/common"
)

type Packet struct {
	ed2kPacket.Packet
}

func NewPacket(
	opCode ed2kCommon.Operation,
) *Packet {
	p := &Packet{
		Packet: *ed2kPacket.NewPacket(common.ProtocolKadUDP, 2),
	}
	p.SetCommand(byte(opCode))
	return p
}

func NewFixedSizePacket(
	opCode ed2kCommon.Operation,
	size int,
) *Packet {
	fmt.Println("kad.NewFixedSizePacket", opCode, size)
	p := &Packet{
		Packet: *ed2kPacket.NewPacket(common.ProtocolKadUDP, size+2),
	}
	p.SetCommand(byte(opCode))
	return p
}

func (packet *Packet) SetCommand(command byte) {
	if packet.GetProtocol() == ed2kCommon.ProtocolEd2kServerUDP || packet.GetProtocol() == common.ProtocolKadUDP {
		packet.SetByte(1, command)
	} else {
		packet.SetByte(5, command)
	}
}

func (packet *Packet) GetCommand() byte {
	if packet.GetProtocol() == ed2kCommon.ProtocolEd2kServerUDP || packet.GetProtocol() == common.ProtocolKadUDP {
		return packet.GetByte(1)
	} else {
		return packet.GetByte(5)
	}
}
