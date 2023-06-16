package udp

import (
	"encoding/binary"
	"errors"
)

type Packet interface {
	GetData() []byte
	GetSize() int

	SetByte(position int, value byte) error
	GetByte(position int) byte
	SafeGetByte(position int) (byte, error)

	AppendUInt16(value uint16) error
	AppendInt(value int) error
}

type RawPacket struct {
	data      []byte
	dataSeek  int
	fixedSize int
}

var _ Packet = &RawPacket{}

// NewPacket creates a new packet
func NewPacket() Packet {
	return NewRawPacket()
}

func NewRawPacket() *RawPacket {
	return &RawPacket{
		data:      make([]byte, 0),
		dataSeek:  0,
		fixedSize: -1,
	}
}

// NewFixedSizePacket creates a new packet and allocates the specified size
func NewFixedSizePacket(size int) Packet {
	return NewFixedSizeRawPacket(size)
}

func NewFixedSizeRawPacket(size int) *RawPacket {
	return &RawPacket{
		data:      make([]byte, size),
		dataSeek:  0,
		fixedSize: size,
	}
}

func (packet *RawPacket) GetData() []byte {
	return packet.data
}

func (packet *RawPacket) GetSize() int {
	if packet.fixedSize >= 0 {
		return packet.fixedSize
	} else {
		return len(packet.data)
	}
}

func (packet *RawPacket) SetByte(position int, value byte) error {
	if packet.fixedSize >= 0 && position >= packet.fixedSize {
		return errors.New("out of bounds")
	}
	if position > len(packet.data) {
		extraSize := position - packet.dataSeek + 1
		packet.data = append(packet.data, make([]byte, extraSize)...)
	}
	packet.data[position] = value
	packet.dataSeek = position + 1
	return nil
}

func (packet *RawPacket) GetByte(position int) byte {
	if position >= len(packet.data) {
		return 0
	}
	return packet.data[position]
}

func (packet *RawPacket) SafeGetByte(position int) (byte, error) {
	if position >= len(packet.data) {
		return 0, errors.New("out of bounds")
	}
	return packet.data[position], nil
}

func (packet *RawPacket) AppendUInt16(value uint16) error {
	if packet.fixedSize >= 0 && len(packet.data)+2 > packet.fixedSize {
		return errors.New("packet is full")
	}

	packet.data = binary.LittleEndian.AppendUint16(packet.data, value)
	packet.dataSeek = len(packet.data)
	return nil
}

func (packet *RawPacket) AppendInt(value int) error {
	if packet.fixedSize >= 0 && len(packet.data)+4 > packet.fixedSize {
		return errors.New("packet is full")
	}

	packet.data = binary.LittleEndian.AppendUint32(packet.data, uint32(value))
	packet.dataSeek = len(packet.data)
	return nil
}
