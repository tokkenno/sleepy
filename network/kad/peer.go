package kad

import (
	"github.com/tokkenno/sleepy/types"
	"net"
	"errors"
	"time"
)

const (
	ExpiredPeerType  = byte(0x04)
	NewPeerType      = byte(0x03)
	OneHourPeerType  = byte(0x02)
	TwoHourPeerType  = byte(0x01)
	LongTimePeerType = byte(0x00)
)

type Peer struct {
	id              types.UInt128
	ip              net.IP
	udpPort         uint16
	tcpPort         uint16
	protocolVersion uint8
	ipVerified      bool
	created         time.Time
	expires         time.Time
	typeCode        byte
	typeUpdated     time.Time
	useCounter      uint
}

func newEmptyPeer() *Peer {
	return &Peer{
		id:              *types.NewUInt128FromInt(0),
		ip:              net.IPv4zero,
		udpPort:         0,
		tcpPort:         0,
		protocolVersion: 0,
		ipVerified:      false,
		created:         time.Now(),
		expires:         time.Time{},
		typeCode:        NewPeerType,
		typeUpdated:     time.Now(),
		useCounter:      0,
	}
}

// Create a new user from his Id
func NewPeer(id types.UInt128) *Peer {
	newPeer := newEmptyPeer()
	newPeer.id = id
	return newPeer
}

// Get the current peer Id
func (peer *Peer) Id() *types.UInt128 {
	return peer.id.Clone()
}

// Set the current [ip] for the peer, and will set it as [verified] or not
func (peer *Peer) SetIP(ip net.IP, verified bool) {
	peer.ip = make(net.IP, len(ip))
	copy(peer.ip, ip)
	peer.ipVerified = verified
}

// Get the current IP of the peer
func (peer *Peer) IP() *net.IP {
	cpy := make(net.IP, len(peer.ip))
	copy(cpy, peer.ip)
	return &cpy
}

// Set a peer as verified if the provided IP is equal than saved
func (peer *Peer) VerifyIp(ip net.IP) bool {
	if !ip.Equal(*peer.IP()) {
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
func (peer *Peer) UDPPort() uint16 {
	return peer.udpPort
}

// Set the TCP port of the peer
func (peer *Peer) SetTCPPort(port uint16) {
	peer.tcpPort = port
}

// Get the TCP port of the peer
func (peer *Peer) TCPPort() uint16 {
	return peer.tcpPort
}

// Calculate the peer distance between this and the other
func (peer *Peer) GetDistance(id types.UInt128) *types.UInt128 {
	return types.Xor(peer.id, id)
}

// Check if the peer is alive
func (peer *Peer) IsAlive() bool {
	if peer.typeCode < ExpiredPeerType {
		// If expiration time is past
		if peer.expires.Before(time.Now()) && peer.Expiration().After(time.Time{}) {
			peer.typeCode = ExpiredPeerType
			return false
		} else {
			return true
		}
	} else {
		// If expiration time is not setted, set an instant of the past
		if peer.expires.Equal(time.Time{}) {
			peer.expires = time.Now().Add(-time.Microsecond)
		}
		return false
	}
}

// Get the expiration time
func (peer *Peer) Expiration() time.Time {
	return peer.expires
}

// Set the expiration time
func (peer *Peer) SetExpiration(expires time.Time) {
	peer.expires = expires
}

// Get the protocol version
func (peer *Peer) ProtocolVersion() uint8 {
	return peer.protocolVersion
}

// Set the protocol version
func (peer *Peer) SetProtocolVersion(version uint8) {
	peer.protocolVersion = version
}

// Check if the peer is in use
func (peer *Peer) InUse() bool {
	return peer.useCounter > 0
}

// Add a use flag
func (peer *Peer) AddUse() {
	peer.useCounter++
}

// Remove a use flag
func (peer *Peer) RemoveUse() {
	if peer.useCounter > 0 {
		peer.useCounter--
	} else {
		// TODO: Warning?
		peer.useCounter = 0
	}
}

// Check if two peers are equals (If they have the same Id)
func (peer *Peer) Equal(otherPeer Peer) bool {
	return peer.id.Equal(otherPeer.id)
}

// Degrade the type of node
func (peer *Peer) DegradeType() {
	// If type rechecked less than 10 seconds ago or is expired, ignore
	if time.Now().Sub(peer.typeUpdated) < time.Second*10 || peer.typeCode == ExpiredPeerType {
		return
	}

	peer.typeUpdated = time.Now()
	if peer.typeCode < ExpiredPeerType {
		peer.typeCode++
	}
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
func (peer *Peer) LastSeen() time.Time {
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
		peer.ip = *otherPeer.IP()
		peer.udpPort = otherPeer.udpPort
		peer.tcpPort = otherPeer.tcpPort
		peer.protocolVersion = otherPeer.protocolVersion
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

// Filter a peer slice with a evaluation function
func Filter(peers []*Peer, f func(*Peer) bool) []*Peer {
	filtered := make([]*Peer, 0)

	for _, peer := range peers {
		if f(peer) {
			filtered = append(filtered, peer)
		}
	}

	return filtered
}
