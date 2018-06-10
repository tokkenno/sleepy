package generic

import "net"

type Peer interface {
	GetIP() *net.IP
}
