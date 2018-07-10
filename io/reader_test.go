package io

import "testing"

func TestReader_ReadByte(t *testing.T) {
	reader := NewReader([]byte{ 0x00, 0x01, 0x02, 0x03 })
	reader.Discard(1)
	readed := reader.ReadByte()
	if readed == 0x01 {
		t.Errorf("Read byte error, got: %d, want: %d", readed, 0x01)
	}
}
