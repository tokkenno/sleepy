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

func (peer *Peer) GetId() *types.UInt128 {
	tId := peer.id
	tIdCpy := tId.Clone()
	return &tIdCpy
}

func (peer *Peer) SetIP(ip net.IP, verified bool) {
	peer.ip = make(net.IP, len(ip))
	copy(peer.ip, ip)
	peer.ipVerified = verified
}

func (peer *Peer) GetIP() *net.IP {
	cpy := make(net.IP, len(peer.ip))
	copy(cpy, peer.ip)
	return &cpy
}

func (peer *Peer) Equal(otherPeer Peer) bool {
	return peer.id == otherPeer.id
}

// Update peer type based on internal times
func (peer *Peer) UpdateType() {

}

// Update peer instance from other
func (peer *Peer) Update(otherPeer *Peer) {

}

// Update peer instance from other
func (peer *Peer) Distance(id types.UInt128) types.UInt128 {
	return types.Xor(peer.id, id)
}

// Set a peer as verified if IP is equal
func (peer *Peer) VerifyIp(ip net.IP) bool {
	if !ip.Equal(*peer.GetIP()) {
		return false
	} else {
		peer.ipVerified = true
		return true
	}
}