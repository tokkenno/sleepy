package io

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/reimashi/sleepy/types"
	"testing"
)

type Reader struct {
	reader *bufio.Reader
	errors []error
}

func NewReader(data []byte) *Reader {
	return &Reader{
		bufio.NewReader(bytes.NewReader(data)),
		make([]error, 0),
	}
}

func (this *Reader) Correct() bool {
	return len(this.errors) == 0
}

func (this *Reader) Discard (count int) {
	discarded, err := this.reader.Discard(count)
	if err != nil {
		this.errors = append(this.errors, err)
	} else if discarded != count {
		this.errors = append(this.errors, errors.New("discarded size mismatch"))
	}
}

func (this *Reader) ReadByte() byte {
	buffer, err := this.reader.ReadByte()
	if (err != nil) {
		this.errors = append(this.errors, err)
		return 0x00
	} else {
		return buffer
	}
}

func (this *Reader) ReadBytes(size uint) []byte {
	buffer := make([]byte, size)
	count, err := this.reader.Read(buffer)
	if err != nil {
		this.errors = append(this.errors, err)
		return make([]byte, size)
	} else if uint(count) != size {
		this.errors = append(this.errors, errors.New("size mismatch"))
		return make([]byte, size)
	} else {
		return buffer
	}
}

func (this *Reader) ReadUInt32() uint32 {
	buffer := this.ReadBytes(4)
	return uint32(buffer[0]) + uint32(buffer[1])<<8 + uint32(buffer[2])<<16 + uint32(buffer[3])<<24
}

func (this *Reader) ReadInt32() int32 {
	return int32(this.ReadUInt32())
}

func (this *Reader) ReadInt() int {
	return int(this.ReadInt32())
}

func (this *Reader) ReadUInt128() (*types.UInt128, error) {
	return types.NewUInt128FromByteArray(this.ReadBytes(16))
}

func TestReadByte (t *testing.T) {
	reader := NewReader([]byte{ 0x00, 0x01, 0x02, 0x03 })
	reader.Discard(1)
	readed := reader.ReadByte()
	if readed == 0x01 {
		t.Errorf("Read byte error, got: %d, want: %d", readed, 0x01)
	}
}