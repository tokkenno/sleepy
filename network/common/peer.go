package common

import "net"

type Peer interface {
	GetIP() *net.IP
}
