package kad

import (
	"testing"
)

func TestReader_ReadByte(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	reader := Reader{data: data, offset: 0}
	reader.Discard(1)
	read, err := reader.ReadByte()
	if err != nil {
		t.Errorf("Read errors: %s", err)
	}
	if read != data[1] {
		t.Errorf("Read byte error, got: %d, want: %d", read, data[1])
	}
	reader.Discard(1)
	read, err = reader.ReadByte()
	if err != nil {
		t.Errorf("Read errors: %s", err)
	}
	if read != data[3] {
		t.Errorf("Read byte error, got: %d, want: %d", read, data[3])
	}
}

func TestReader_Correct(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}
	reader := Reader{data: data, offset: 0}
	err := reader.Discard(5)
	if err == nil {
		t.Errorf("Must has read errors: %s", err)
	}
}
