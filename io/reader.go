package io

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/tokkenno/sleepy/types"
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

func (reader *Reader) GetErrors() []error {
	return reader.errors
}

func (reader *Reader) Correct() bool {
	return len(reader.errors) == 0
}

func (reader *Reader) Discard(count int) {
	discarded, err := reader.reader.Discard(count)
	if err != nil {
		reader.errors = append(reader.errors, err)
	} else if discarded != count {
		reader.errors = append(reader.errors, errors.New("discarded size mismatch"))
	}
}

func (reader *Reader) ReadByte() byte {
	buffer, err := reader.reader.ReadByte()
	if err != nil {
		reader.errors = append(reader.errors, err)
		return 0x00
	} else {
		return buffer
	}
}

func (reader *Reader) ReadBytes(size uint) []byte {
	buffer := make([]byte, size)
	count, err := reader.reader.Read(buffer)
	if err != nil {
		reader.errors = append(reader.errors, err)
		return make([]byte, size)
	} else if uint(count) != size {
		reader.errors = append(reader.errors, errors.New("size mismatch"))
		return make([]byte, size)
	} else {
		return buffer
	}
}

func (reader *Reader) ReadUInt32() uint32 {
	buffer := reader.ReadBytes(4)
	return uint32(buffer[0]) + uint32(buffer[1])<<8 + uint32(buffer[2])<<16 + uint32(buffer[3])<<24
}

func (reader *Reader) ReadUInt16() uint16 {
	buffer := reader.ReadBytes(2)
	return uint16(buffer[0]) + uint16(buffer[1])<<8
}

func (reader *Reader) ReadInt32() int32 {
	return int32(reader.ReadUInt32())
}

func (reader *Reader) ReadInt() int {
	return int(reader.ReadInt32())
}

func (reader *Reader) ReadUInt128() (*types.UInt128, error) {
	return types.NewUInt128FromByteArray(reader.ReadBytes(16))
}

func (reader *Reader) ReadString(txtSize uint) string {
	return string(reader.ReadBytes(uint(txtSize)))
}

func (reader *Reader) ReadTags() map[interface{}]interface{} {
	tags := make(map[interface{}]interface{})

	tagCount := reader.ReadUInt32()
	for ind := uint32(0); ind < tagCount; ind++ {
		tagType := reader.ReadByte()

		var key interface{}
		keySize := reader.ReadUInt16()
		if keySize == 1 {
			key = uint8(reader.ReadByte())
		} else {
			key = reader.ReadString(uint(keySize))
		}

		switch tagType {
		case 0x02:
			tags[key] = reader.ReadString(uint(reader.ReadUInt16()))
			break
		case 0x03:
			tags[key] = reader.ReadInt32()
			break
		case 0x04:
			tags[key] = reader.ReadInt32() // FIXME: Float
			break
		default:
			reader.errors = append(reader.errors, errors.New("Tag no reconocido: %d\n"))
		}
	}

	return tags
}
