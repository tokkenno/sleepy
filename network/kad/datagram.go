package kad

import (
	"net"
	"bytes"
)

type Datagram struct {
	typeCode byte
	data []byte
}

type InDatagram struct {
	Datagram
	reader *bytes.Reader
	from net.Addr
}

type KadInDatagram struct {
	*InDatagram
	command byte
}