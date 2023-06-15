package kad

import (
	"errors"
	"sleepy/types"
)

type Reader struct {
	data   []byte
	offset int64
}

func (reader *Reader) Read(buffer []byte) (n int, err error) {
	if reader.offset < int64(len(reader.data)) {
		bytesRead := copy(buffer, reader.data[reader.offset:])
		reader.offset = reader.offset + int64(bytesRead)
		return bytesRead, nil
	} else {
		return 0, errors.New("out of bounds")
	}
}

func (reader *Reader) Seek(offset int64, whence int) (int64, error) {
	if offset >= int64(len(reader.data)) || offset < 0 {
		return 0, errors.New("out of bounds")
	} else {
		reader.offset = offset
		return offset, nil
	}
}

func (reader *Reader) Close() error {
	return nil
}

func (reader *Reader) Discard(size int64) error {
	_, err := reader.Seek(reader.offset+size, 0)
	return err
}

func (reader *Reader) ReadByte() (byte, error) {
	buffer, err := reader.ReadBytes(1)
	if err != nil {
		return 0x00, err
	} else {
		return buffer[0], nil
	}
}

func (reader *Reader) ReadBytes(size uint) ([]byte, error) {
	buffer := make([]byte, size)
	count, err := reader.Read(buffer)
	if err != nil {
		return nil, err
	} else if uint(count) != size {
		return buffer, errors.New("size mismatch")
	} else {
		return buffer, nil
	}
}

func (reader *Reader) ReadUInt32() (uint32, error) {
	buffer, err := reader.ReadBytes(4)
	if err != nil {
		return 0, err
	} else {
		return uint32(buffer[0]) + uint32(buffer[1])<<8 + uint32(buffer[2])<<16 + uint32(buffer[3])<<24, nil
	}
}

func (reader *Reader) ReadUInt16() (uint16, error) {
	buffer, err := reader.ReadBytes(2)
	if err != nil {
		return 0, err
	} else {
		return uint16(buffer[0]) + uint16(buffer[1])<<8, nil
	}
}

func (reader *Reader) ReadInt32() (int32, error) {
	value, err := reader.ReadUInt32()
	return int32(value), err
}

func (reader *Reader) ReadInt() (int, error) {
	value, err := reader.ReadUInt32()
	return int(value), err
}

func (reader *Reader) ReadUInt128() (types.UInt128, error) {
	buffer, err := reader.ReadBytes(16)
	if err != nil {
		return nil, err
	} else {
		return types.NewUInt128FromByteArray(buffer)
	}
}

func (reader *Reader) ReadString(txtSize uint) (string, error) {
	buffer, err := reader.ReadBytes(txtSize)
	if err != nil {
		return "", err
	} else {
		return string(buffer), nil
	}
}

func (reader *Reader) ReadTags() (map[interface{}]interface{}, error) {
	tags := make(map[interface{}]interface{})

	tagCount, err := reader.ReadUInt32()

	if err != nil {
		return nil, err
	}

	for ind := uint32(0); ind < tagCount; ind++ {
		tagType, err := reader.ReadByte()

		if err != nil {
			return nil, err
		}

		var key interface{}
		keySize, err := reader.ReadUInt16()

		if err != nil {
			return nil, err
		} else if keySize == 1 {
			byte, err := reader.ReadByte()

			if err != nil {
				return nil, err
			} else {
				key = uint8(byte)
			}
		} else {
			key, err = reader.ReadString(uint(keySize))

			if err != nil {
				return nil, err
			}
		}

		switch tagType {
		case 0x02:
			valueSize, err := reader.ReadUInt16()
			if err != nil {
				return nil, err
			}
			tags[key], err = reader.ReadString(uint(valueSize))
			if err != nil {
				return nil, err
			}
			break
		case 0x03:
			tags[key], err = reader.ReadInt32()
			if err != nil {
				return nil, err
			}
			break
		case 0x04:
			tags[key], err = reader.ReadInt32() // FIXME: Float
			if err != nil {
				return nil, err
			}
			break
		default:
			return nil, errors.New("unknown tag type")
		}
	}

	return tags, nil
}
