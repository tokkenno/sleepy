package kad

import (
	"net"
	"com/github/reimashi/sleepy/io"
)

type Datagram struct {
	protocolCode byte
	data *[]byte
}

type InDatagram struct {
	Datagram
	reader *io.Reader
	from net.Addr
}

type KadInDatagram struct {
	*InDatagram
	command byte
}