package udp

import "encoding/binary"

type Packet interface {
	GetData() []byte
}

type RawPacket struct {
	data []byte
}

var _ Packet = &RawPacket{}

func NewPacket(size int) *RawPacket {
	return &RawPacket{
		data: make([]byte, 0, size),
	}
}

func (packet *RawPacket) GetData() []byte {
	return packet.data
}

func (packet *RawPacket) Size() int {
	return len(packet.data)
}

func (packet *RawPacket) init() {
	if packet.data == nil {
		packet.data = make([]byte, 0)
	}
}

func (packet *RawPacket) SetByte(position int, value byte) {
	packet.init()
	if position > len(packet.data) {
		packet.data = append(packet.data, make([]byte, position-len(packet.data))...)
	}
	packet.data[position] = value
}

func (packet *RawPacket) GetByte(position int) byte {
	if position >= len(packet.data) {
		return 0
	}
	return packet.data[position]
}

func (packet *RawPacket) AppendUInt16(value uint16) {
	packet.init()
	binary.LittleEndian.AppendUint16(packet.data, value)
}

func (packet *RawPacket) AppendInt(value int) {
	packet.init()
	binary.LittleEndian.AppendUint32(packet.data, uint32(value))
}
