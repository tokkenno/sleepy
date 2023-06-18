package types

import (
	"errors"
	"net"
	"sleepy/types"
	"time"
)

const (
	ExpiredPeerType  = byte(0x04)
	NewPeerType      = byte(0x03)
	OneHourPeerType  = byte(0x02)
	TwoHourPeerType  = byte(0x01)
	LongTimePeerType = byte(0x00)
)

type Peer interface {
	Equal(other Peer) bool
	GetID() types.UInt128
	// SetIP set the current [ip] for the peer, and will set it as [verified] or not
	SetIP(ip net.IP, verified bool)
	// GetIP get the current IP of the peer
	GetIP() net.IP
	// VerifyIp set a peer as verified if the provided IP is equal than saved
	VerifyIp(ip net.IP) bool
	IsIPVerified() bool
	IsAlive() bool
	// InUse Check if the peer is in use
	InUse() bool
	GetUDPPort() uint16
	SetUDPPort(port uint16)
	GetTCPPort() uint16
	SetTCPPort(port uint16)
	GetProtocolVersion() uint8
	SetProtocolVersion(version uint8)
	GetCreatedAt() time.Time
	GetExpiresAt() time.Time
	// SetExpiration set the expiration time
	SetExpiration(ea time.Time)
	// GetDistance calculates the peer distance between this and the other
	GetDistance(id types.UInt128) types.UInt128
	GetTypeCode() byte
	GetTypeUpdatedAt() time.Time
	// DegradeType degrade the type of node
	DegradeType()
	UpdateType()
	// UpdateFrom update peer instance from other
	UpdateFrom(otherPeer Peer) error
}

type peerImp struct {
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

func newEmptyPeer() *peerImp {
	return &peerImp{
		id:              types.NewUInt128FromInt(0),
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

// NewPeer create a new user from his id
func NewPeer(id types.UInt128) Peer {
	newPeer := newEmptyPeer()
	newPeer.id = id.Clone()
	return newPeer
}

// Id gets the current peer id
func (peer *peerImp) GetID() types.UInt128 {
	return peer.id.Clone()
}

func (peer *peerImp) SetIP(ip net.IP, verified bool) {
	peer.ip = make(net.IP, len(ip))
	copy(peer.ip, ip)
	peer.ipVerified = verified
}

func (peer *peerImp) GetIP() net.IP {
	cpy := make(net.IP, len(peer.ip))
	copy(cpy, peer.ip)
	return cpy
}

func (peer *peerImp) VerifyIp(ip net.IP) bool {
	if !ip.Equal(peer.GetIP()) {
		peer.ipVerified = false
		return false
	} else {
		peer.ipVerified = true
		return true
	}
}

// Check if the current IP is verified
func (peer *peerImp) IsIPVerified() bool {
	return peer.ipVerified
}

// Set the UDP port of the peer
func (peer *peerImp) SetUDPPort(port uint16) {
	peer.udpPort = port
}

// Get the UDP port of the peer
func (peer *peerImp) GetUDPPort() uint16 {
	return peer.udpPort
}

// Set the TCP port of the peer
func (peer *peerImp) SetTCPPort(port uint16) {
	peer.tcpPort = port
}

// Get the TCP port of the peer
func (peer *peerImp) GetTCPPort() uint16 {
	return peer.tcpPort
}

func (peer *peerImp) GetDistance(id types.UInt128) types.UInt128 {
	return types.Xor(peer.id, id)
}

// Check if the peer is alive
func (peer *peerImp) IsAlive() bool {
	if peer.typeCode < ExpiredPeerType {
		// If expiration time is past
		if peer.expires.Before(time.Now()) && peer.GetExpiresAt().After(time.Time{}) {
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

func (peer *peerImp) GetCreatedAt() time.Time {
	return peer.created
}

func (peer *peerImp) GetExpiresAt() time.Time {
	return peer.expires
}

func (peer *peerImp) SetExpiration(expires time.Time) {
	peer.expires = expires
}

func (peer *peerImp) GetTypeCode() byte {
	return peer.typeCode
}

func (peer *peerImp) GetTypeUpdatedAt() time.Time {
	return peer.typeUpdated
}

// Get the protocol version
func (peer *peerImp) GetProtocolVersion() uint8 {
	return peer.protocolVersion
}

// Set the protocol version
func (peer *peerImp) SetProtocolVersion(version uint8) {
	peer.protocolVersion = version
}

func (peer *peerImp) InUse() bool {
	return peer.useCounter > 0
}

// Add a use flag
func (peer *peerImp) AddUse() {
	peer.useCounter++
}

// Remove a use flag
func (peer *peerImp) RemoveUse() {
	if peer.useCounter > 0 {
		peer.useCounter--
	} else {
		// TODO: Warning?
		peer.useCounter = 0
	}
}

// Check if two peers are equals (If they have the same Id)
func (peer *peerImp) Equal(otherPeer Peer) bool {
	return peer.id.Equal(otherPeer.GetID())
}

func (peer *peerImp) DegradeType() {
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
func (peer *peerImp) UpdateType() {
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
func (peer *peerImp) LastSeen() time.Time {
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

func (peer *peerImp) UpdateFrom(otherPeer Peer) error {
	if peer.Equal(otherPeer) {
		peer.ip = otherPeer.GetIP()
		peer.udpPort = otherPeer.GetUDPPort()
		peer.tcpPort = otherPeer.GetTCPPort()
		peer.protocolVersion = otherPeer.GetProtocolVersion()
		peer.ipVerified = otherPeer.IsIPVerified()
		peer.created = otherPeer.GetCreatedAt()
		peer.expires = otherPeer.GetExpiresAt()
		peer.typeCode = otherPeer.GetTypeCode()
		peer.typeUpdated = otherPeer.GetTypeUpdatedAt()
		return nil
	} else {
		return errors.New("the peer information only can be updated with the information of other peer with the same id")
	}
}

// Filter a peer slice with an evaluation function
func Filter(peers []Peer, f func(Peer) bool) []Peer {
	filtered := make([]Peer, 0)

	for _, peer := range peers {
		if f(peer) {
			filtered = append(filtered, peer)
		}
	}

	return filtered
}
