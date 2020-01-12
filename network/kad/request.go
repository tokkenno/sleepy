package kad

import (
	"context"
	"net"
)

type Request struct {
	body Reader
	ctx  context.Context
}

type UDPRequest struct {
	Request
	from *net.UDPAddr
}