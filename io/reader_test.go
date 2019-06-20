package io

import "testing"

func TestReader_ReadByte(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	reader := NewReader(data)
	reader.Discard(1)
	read := reader.ReadByte()
	if read != data[1] {
		t.Errorf("Read byte error, got: %d, want: %d", read, data[1])
	}
	reader.Discard(1)
	read = reader.ReadByte()
	if read != data[3] {
		t.Errorf("Read byte error, got: %d, want: %d", read, data[3])
	}
	if !reader.Correct() {
		t.Errorf("Read errors: %s", reader.errors[0].Error())
	}
}

func TestReader_Correct(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	reader := NewReader(data)
	reader.Discard(5)
	if reader.Correct() {
		t.Errorf("No errors. EOF error expected.")
	}
}
