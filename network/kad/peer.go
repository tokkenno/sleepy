package kad

import (
	"github.com/reimashi/sleepy/types"
	"net"
)

type Peer struct {
	Id types.UInt128
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
func (this *Peer) Distance(id types.UInt128) types.UInt128 {
	return types.Xor(this.Id, id)
}