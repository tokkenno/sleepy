package kad

import (
	"context"
	"net"
	"sleepy/types"
)

type Request struct {
	body Reader
	ctx  context.Context
}

type UDPRequest struct {
	Request
	from *net.UDPAddr
}

type Bootstrap1Request struct {
	Request
	ClientID types.UInt128
	Address  net.IP
	UDPPort  uint16
	TCPPort  uint16
}
