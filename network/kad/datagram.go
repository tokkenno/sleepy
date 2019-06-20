package kad

import (
	"github.com/tokkenno/sleepy/io"
	"net"
)

type Datagram struct {
	protocolCode byte
	data         *[]byte
}

type InDatagram struct {
	Datagram
	reader *io.Reader
	from   net.Addr
}

type KadInDatagram struct {
	*InDatagram
	command byte
}
