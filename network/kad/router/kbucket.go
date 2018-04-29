package router

import (
	"errors"
	"com/github/reimashi/sleepy/network/kad"
	"com/github/reimashi/sleepy/types/uint128"
)

const (
	maxSize = 16		// Max number of peers in each k-bucket
	maxPeersPerIP = 3 	// Number of peers permitted from same public IP
)

// K-bucket is a queue of k peers ordered by TTL
type KBucket struct {
	peers []kad.Peer
}

// Count the number of peers on this k-bucket
func (this *KBucket) CountPeers() int {
	return len(this.peers)
}

// Check peer data and append to the end of k-bucket if correct
func (this *KBucket) AddPeer(newPeer *kad.Peer) error {
	if newPeer == nil { return errors.New("KBucket only can storage not null peers") }

	sameNetwork := 0
	for _, peer := range this.peers {
		if peer.Equal(newPeer) { return errors.New("KBucket already contains this peer") }
		if peer.GetIP().Equal(newPeer.GetIP()) { sameNetwork++ }
	}

	if (len(this.peers) >= maxSize) {
		return errors.New("The current KBucket is full")
	}

	if (sameNetwork >= maxPeersPerIP) {
		return errors.New("Many peers for the current IP")
	}

	this.peers = append(this.peers, *newPeer)
	// TODO: Adjust global tracking
	return nil
}

// Get a peer from his id
func (this *KBucket) GetPeer(id *uint128.UInt128) (*kad.Peer, error) {
	for _, peer := range this.peers {
		if peer.Id.Equal(id) {
			return &peer, nil
		}
	}
	return nil, errors.New("KBucket not constains a peer with this id")
}

// Push the peer to the end of bucket
func (this *KBucket) pushToEnd(peer *kad.Peer) error {
	for position, currPeer := range this.peers {
		if currPeer.Equal(peer) {
			this.peers = append(append(this.peers[0:position], this.peers[position+1:len(this.peers)]...), *peer)
			return nil
		}
	}
	return errors.New("KBucket not constains this peer")
}

// Update last viewed status of peer, recalculate and set TTL
func (this *KBucket) SetAlive(peer *kad.Peer) error {
	inPeer, err := this.GetPeer(&peer.Id)
	if (err != nil) { return err }

	inPeer.UpdateType()
	return this.pushToEnd(inPeer)
}