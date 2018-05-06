package router

import (
	"errors"
	"com/github/reimashi/sleepy/network/kad"
	"com/github/reimashi/sleepy/types"
	"net"
	"sync"
	"math/rand"
)

const (
	maxSize = 16		// Max number of peers in each k-bucket
	maxPeersPerIP = 3 	// Number of peers permitted from same public IP
)

// K-bucket is a queue of k peers ordered by TTL
type KBucket struct {
	peers []kad.Peer
	peersAccess sync.Mutex
}

// Count the number of peers on this k-bucket
func (this *KBucket) Count() int {
	return len(this.peers)
}

// Check peer data and append to the end of k-bucket if correct
func (this *KBucket) AddPeer(newPeer *kad.Peer) error {
	if newPeer == nil { return errors.New("KBucket only can storage not null peers") }

	this.peersAccess.Lock()

	sameNetwork := 0
	for _, peer := range this.peers {
		if peer.Equal(newPeer) {
			this.peersAccess.Unlock()
			return errors.New("KBucket already contains this peer")
		}
		if peer.GetIP().Equal(newPeer.GetIP()) { sameNetwork++ }
	}

	if (len(this.peers) >= maxSize) {
		this.peersAccess.Unlock()
		return errors.New("The current KBucket is full")
	}

	if (sameNetwork >= maxPeersPerIP) {
		this.peersAccess.Unlock()
		return errors.New("Many peers for the current IP")
	}

	this.peers = append(this.peers, *newPeer)
	this.peersAccess.Unlock()
	// TODO: Adjust global tracking
	return nil
}

// Remove a peer from the bucket
func (this *KBucket) RemovePeer(peer *kad.Peer) error {
	this.peersAccess.Lock()

	for index, peertr := range this.peers {
		if &peertr == peer {
			this.peers = append(this.peers[0:index], this.peers[index+1:len(this.peers)]...)
			this.peersAccess.Unlock()
			return nil
		}
	}

	this.peersAccess.Unlock()
	return errors.New("KBucket not constains a peer with this id")
}

// Get a peer from his id
func (this *KBucket) GetPeer(id types.UInt128) (*kad.Peer, error) {
	this.peersAccess.Lock()

	for _, peer := range this.peers {
		if peer.Id.Equal(id) {
			this.peersAccess.Unlock()
			return &peer, nil
		}
	}

	this.peersAccess.Unlock()
	return nil, errors.New("KBucket not constains a peer with this id")
}

// Get a peer from his ip
func (this *KBucket) GetPeerByIp(ip net.IP) (*kad.Peer, error) {
	this.peersAccess.Lock()

	for _, peer := range this.peers {
		if peer.GetIP().Equal(ip) {
			this.peersAccess.Unlock()
			return &peer, nil
		}
	}

	this.peersAccess.Unlock()
	return nil, errors.New("KBucket not constains a peer with this ip")
}

func (this *KBucket) GetRandomPeer() (*kad.Peer, error) {
	this.peersAccess.Lock()

	if (len(this.peers) > 0) {
		peer := &this.peers[rand.Intn(len(this.peers))]
		this.peersAccess.Unlock()
		return peer, nil
	} else {
		this.peersAccess.Unlock()
		return nil, errors.New("KBucket not constains any peer")
	}
}

func (this *KBucket) GetPeers() []kad.Peer {
	return this.peers
}

func (this *KBucket) IsFull() bool {
	return this.Count() == maxSize
}

// Push the peer to the end of bucket (only if already exists)
func (this *KBucket) pushToEnd(peer *kad.Peer) error {
	this.peersAccess.Lock()

	for position, currPeer := range this.peers {
		if &currPeer == peer {
			this.peers = append(append(this.peers[0:position], this.peers[position+1:len(this.peers)]...), *peer)
			this.peersAccess.Unlock()
			return nil
		}
	}

	this.peersAccess.Unlock()
	return errors.New("KBucket not constains this peer")
}

// Update last viewed status of peer, recalculate and set TTL
func (this *KBucket) SetAlive(id types.UInt128) error {
	inPeer, err := this.GetPeer(id)
	if (err != nil) { return err }

	inPeer.UpdateType()
	return this.pushToEnd(inPeer)
}