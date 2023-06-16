package udp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFixedSizePacket(t *testing.T) {
	packet := NewFixedSizePacket(2)
	// Append must fail because packet is fixed size
	assert.Error(t, packet.AppendUInt16(1))
	assert.Error(t, packet.SetByte(2, 0))
	assert.NoError(t, packet.SetByte(1, 0))
}

func TestRawPacket_AppendInt(t *testing.T) {
	packet := NewPacket()
	assert.NoError(t, packet.AppendInt(1))
	assert.NoError(t, packet.AppendInt(2))
	assert.EqualValues(t, []byte{1, 0, 0, 0, 2, 0, 0, 0}, packet.GetData())
}

func TestRawPacket_AppendUInt16(t *testing.T) {
	packet := NewPacket()
	assert.NoError(t, packet.AppendUInt16(1))
	assert.NoError(t, packet.AppendUInt16(2))
	assert.EqualValues(t, []byte{1, 0, 2, 0}, packet.GetData())
}

func TestRawPacket_GetByte(t *testing.T) {
	packet := NewPacket()
	assert.NoError(t, packet.AppendUInt16(1))
	assert.NoError(t, packet.AppendUInt16(2))
	assert.EqualValues(t, byte(1), packet.GetByte(0))
	assert.EqualValues(t, byte(2), packet.GetByte(2))
	assert.Equal(t, 4, packet.GetSize())
}

func TestRawPacket_SetByte(t *testing.T) {
	packet := NewPacket()
	assert.NoError(t, packet.SetByte(1, 3))
	assert.NoError(t, packet.SetByte(3, 4))
	assert.EqualValues(t, []byte{0, 3, 0, 4}, packet.GetData())
	assert.Equal(t, 4, packet.GetSize())
}

func TestRawPacket_GetSize(t *testing.T) {
	packet := NewPacket()
	assert.NoError(t, packet.AppendUInt16(1))
	assert.NoError(t, packet.AppendUInt16(2))
	assert.EqualValues(t, 4, packet.GetSize())
}
