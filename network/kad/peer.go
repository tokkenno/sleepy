package kad

import (
	"com/github/reimashi/sleepy/types/uint128"
	"net"
)

type Peer struct {
	Id uint128.UInt128
	ip net.IP
	udpPort uint16
	tcpPort uint16
	uVersion uint8
	IpVerified bool
}

func (this *Peer) Equal(peer *Peer) bool { return this.Id == peer.Id }

func (this *Peer) GetIP() net.IP { return this.ip }

// Update peer type based on internal times
func (this *Peer) UpdateType() {

}

// Update peer instance from other
func (this *Peer) Update(peer *Peer) {

}

// Update peer instance from other
func (this *Peer) Distance(id uint128.UInt128) uint128.UInt128 {
	return uint128.Xor(this.Id, id)
}