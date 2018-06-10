package kad

import (
	"github.com/tokkenno/sleepy/types"
	"net"
	"errors"
	"time"
)

const (
	ElapsedPeerType  = byte(0x04)
	NewPeerType      = byte(0x03)
	OneHourPeerType  = byte(0x02)
	TwoHourPeerType  = byte(0x01)
	LongTimePeerType = byte(0x00)
)

type Peer struct {
	id          types.UInt128
	ip          net.IP
	udpPort     uint16
	tcpPort     uint16
	uVersion    uint8
	ipVerified  bool
	created     time.Time
	expires     time.Time
	typeCode    byte
	typeUpdated time.Time
}

func newEmptyPeer() *Peer {
	return &Peer{
		id:          *types.NewUInt128FromInt(0),
		ip:          net.IPv4zero,
		udpPort:     0,
		tcpPort:     0,
		uVersion:    0,
		ipVerified:  false,
		created:     time.Now(),
		expires:     time.Time{},
		typeCode:    NewPeerType,
		typeUpdated: time.Now(),
	}
}

func NewPeer(id types.UInt128) *Peer {
	newPeer := newEmptyPeer()
	newPeer.id = id
	return newPeer
}

// Get the current peer Id
func (peer *Peer) GetId() *types.UInt128 {
	return peer.id.Clone()
}

// Set the current [ip] for the peer, and will set it as [verified] or not
func (peer *Peer) SetIP(ip net.IP, verified bool) {
	peer.ip = make(net.IP, len(ip))
	copy(peer.ip, ip)
	peer.ipVerified = verified
}

// Get the current IP of the peer
func (peer *Peer) GetIP() *net.IP {
	cpy := make(net.IP, len(peer.ip))
	copy(cpy, peer.ip)
	return &cpy
}

// Set a peer as verified if the provided IP is equal than saved
func (peer *Peer) VerifyIp(ip net.IP) bool {
	if !ip.Equal(*peer.GetIP()) {
		peer.ipVerified = false
		return false
	} else {
		peer.ipVerified = true
		return true
	}
}

// Check if the current IP is verified
func (peer *Peer) IsIpVerified() bool {
	return peer.ipVerified
}

// Set the UDP port of the peer
func (peer *Peer) SetUDPPort(port uint16) {
	peer.udpPort = port
}

// Get the UDP port of the peer
func (peer *Peer) GetUDPPort() uint16 {
	return peer.udpPort
}

// Set the TCP port of the peer
func (peer *Peer) SetTCPPort(port uint16) {
	peer.tcpPort = port
}

// Get the TCP port of the peer
func (peer *Peer) GetTCPPort() uint16 {
	return peer.tcpPort
}

// Calculate the peer distance between this and the other
func (peer *Peer) GetDistance(id types.UInt128) *types.UInt128 {
	return types.Xor(peer.id, id)
}

// Check if two peers are equals (If they have the same Id)
func (peer *Peer) Equal(otherPeer Peer) bool {
	return peer.id.Equal(otherPeer.id)
}

// ¿¿??
func (peer *Peer) CheckingType() {
	// If type updated less than 10 seconds ago or is expired, ignore
	if time.Now().Sub(peer.typeUpdated) < time.Second*10 || peer.typeCode == ElapsedPeerType {
		return
	}

	peer.typeUpdated = time.Now()
	peer.typeCode++
}

// Update peer type based on internal times
func (peer *Peer) UpdateType() {
	hoursOnline := time.Now().Sub(peer.created)

	if hoursOnline > 2*time.Hour {
		peer.typeCode = LongTimePeerType
		peer.expires = time.Now().Add(time.Hour * 2)
	} else if hoursOnline > time.Hour {
		peer.typeCode = TwoHourPeerType
		peer.expires = time.Now().Add(time.Hour + (time.Minute * 30))
	} else {
		peer.typeCode = OneHourPeerType
		peer.expires = time.Now().Add(time.Hour)
	}
}

// Get the time on which peer has been viewed last time
func (peer *Peer) GetLastSeen() time.Time {
	if !peer.expires.Equal(time.Time{}) {
		if peer.typeCode == OneHourPeerType {
			return peer.expires.Add(-time.Hour)
		} else if peer.typeCode == TwoHourPeerType {
			return peer.expires.Add(-time.Hour - (time.Minute * 30))
		} else if peer.typeCode == LongTimePeerType {
			return peer.expires.Add(-time.Hour * 2)
		}
	}
	return time.Time{}
}

// Update peer instance from other
func (peer *Peer) Update(otherPeer *Peer) error {
	if peer.Equal(*otherPeer) {
		peer.ip = *otherPeer.GetIP()
		peer.udpPort = otherPeer.udpPort
		peer.tcpPort = otherPeer.tcpPort
		peer.uVersion = otherPeer.uVersion
		peer.ipVerified = otherPeer.ipVerified
		peer.created = otherPeer.created
		peer.expires = otherPeer.expires
		peer.typeCode = otherPeer.typeCode
		peer.typeUpdated = otherPeer.typeUpdated
		return nil
	} else {
		return errors.New("the peer information only can be updated with the information of other peer with the same id")
	}
}
