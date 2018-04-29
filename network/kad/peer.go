package kad

import (
	"com/github/reimashi/sleepy/types/uint128"
	"net"
)

type Peer struct {
	Id uint128.UInt128
	UDPAddr net.UDPAddr
	TCPAddr net.TCPAddr
	uVersion uint8
	IpVerified bool
}