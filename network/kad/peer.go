package kad

import (
	"github.com/reimashi/sleepy/types"
	"net"
)

type Peer struct {
	id types.UInt128
	ip net.IP
	udpPort uint16
	tcpPort uint16
	uVersion uint8
	ipVerified bool
}

func NewPeer(id types.UInt128) *Peer {
	return &Peer{
		id: id,
		ipVerified: false,
	}
}

func (this *Peer) GetId() *types.UInt128 {
	tId := this.id
	tIdCpy := tId.Clone()
	return &tIdCpy
}

func (this *Peer) SetIP(ip net.IP, verified bool) {
	this.ip = make(net.IP, len(ip))
	copy(this.ip, ip)
	this.ipVerified = verified
}

func (this *Peer) GetIP() *net.IP {
	cpy := make(net.IP, len(this.ip))
	copy(cpy, this.ip)
	return &cpy
}

func (this *Peer) Equal(peer Peer) bool { return this.id == peer.id }

// Update peer type based on internal times
func (this *Peer) UpdateType() {

}

// Update peer instance from other
func (this *Peer) Update(peer *Peer) {

}

// Update peer instance from other
func (this *Peer) Distance(id types.UInt128) types.UInt128 {
	return types.Xor(this.id, id)
}

// Set a peer as verified if IP is equal
func (this *Peer) VerifyIp(ip net.IP) bool {
	if !ip.Equal(*this.GetIP()) {
		return false
	} else {
		this.ipVerified = true
		return true
	}
}