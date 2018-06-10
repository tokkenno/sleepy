package router

import (
	"errors"
	"github.com/tokkenno/sleepy/network/kad"
	"github.com/tokkenno/sleepy/types"
	"net"
	"sync"
	"math/rand"
	"time"
)

const (
	maxSize = 16		// Max number of peers in each k-bucket
	maxPeersPerIP = 3 	// Number of peers permitted from same public IP
)

// K-bucket is a queue of k peers ordered by TTL
type KBucket struct {
	randGen *rand.Rand
	peers []kad.Peer
	peersAccess sync.Mutex
}

// Count the number of peers on this k-bucket
func (bucket *KBucket) CountPeers() int {
	return len(bucket.peers)
}

// Check peer data and append to the end of k-bucket if correct
func (bucket *KBucket) AddPeer(newPeer *kad.Peer) error {
	if newPeer == nil { return errors.New("KBucket only can storage not null peers") }

	bucket.peersAccess.Lock()

	sameNetwork := 0
	for _, peer := range bucket.peers {
		if peer.Equal(*newPeer) {
			bucket.peersAccess.Unlock()
			return errors.New("KBucket already contains the passed peer")
		}
		if peer.GetIP().Equal(*newPeer.GetIP()) { sameNetwork++ }
	}

	if len(bucket.peers) >= maxSize {
		bucket.peersAccess.Unlock()
		return errors.New("the current KBucket is full")
	}

	if sameNetwork >= maxPeersPerIP {
		bucket.peersAccess.Unlock()
		return errors.New("many peers for the current IP")
	}

	bucket.peers = append(bucket.peers, *newPeer)
	bucket.peersAccess.Unlock()
	// TODO: Adjust global tracking
	return nil
}

// Remove a peer from the bucket
func (bucket *KBucket) RemovePeer(peer *kad.Peer) error {
	bucket.peersAccess.Lock()

	for index, peertr := range bucket.peers {
		if peertr.Equal(*peer) {
			bucket.peers = append(bucket.peers[0:index], bucket.peers[index+1:len(bucket.peers)]...)
			bucket.peersAccess.Unlock()
			return nil
		}
	}

	bucket.peersAccess.Unlock()
	return errors.New("KBucket don't contains a peer with the passed id")
}

// Get a peer from his id
func (bucket *KBucket) GetPeer(id types.UInt128) (*kad.Peer, error) {
	bucket.peersAccess.Lock()

	for _, peer := range bucket.peers {
		peerPtr := &peer
		peerId := peerPtr.GetId()

		if peerId.Equal(id) {
			bucket.peersAccess.Unlock()
			return peerPtr, nil
		}
	}

	bucket.peersAccess.Unlock()
	return nil, errors.New("KBucket don't contains a peer with the passed id")
}

// Get a peer from his ip
func (bucket *KBucket) GetPeerByIp(ip net.IP) (*kad.Peer, error) {
	bucket.peersAccess.Lock()

	for _, peer := range bucket.peers {
		if peer.GetIP().Equal(ip) {
			bucket.peersAccess.Unlock()
			return &peer, nil
		}
	}

	bucket.peersAccess.Unlock()
	return nil, errors.New("KBucket don't contains a peer with the passed ip")
}

func (bucket *KBucket) GetRandomPeer() (*kad.Peer, error) {
	if bucket.randGen == nil { bucket.randGen = rand.New(rand.NewSource(time.Now().UnixNano())) }

	bucket.peersAccess.Lock()

	if len(bucket.peers) > 0 {
		peer := &bucket.peers[bucket.randGen.Intn(len(bucket.peers))]
		bucket.peersAccess.Unlock()
		return peer, nil
	} else {
		bucket.peersAccess.Unlock()
		return nil, errors.New("KBucket don't contains any peer")
	}
}

func (bucket *KBucket) ContainsPeer(id types.UInt128) bool {
	peer, err := bucket.GetPeer(id)
	return peer != nil && err != nil
}

func (bucket *KBucket) GetPeers() []kad.Peer {
	return bucket.peers
}

func (bucket *KBucket) IsFull() bool {
	return bucket.CountPeers() == maxSize
}

// Push the peer to the end of bucket (only if already exists)
func (bucket *KBucket) pushToEnd(peer *kad.Peer) error {
	bucket.peersAccess.Lock()

	for position, currPeer := range bucket.peers {
		if peer.Equal(currPeer) {
			bucket.peers = append(append(bucket.peers[0:position], bucket.peers[position+1:len(bucket.peers)]...), *peer)
			bucket.peersAccess.Unlock()
			return nil
		}
	}

	bucket.peersAccess.Unlock()
	return errors.New("KBucket don't contains the passed peer")
}

// Update last viewed status of peer, recalculate and set TTL
func (bucket *KBucket) SetPeerAlive(id types.UInt128) error {
	inPeer, err := bucket.GetPeer(id)
	if err != nil { return err }

	inPeer.UpdateType()
	return bucket.pushToEnd(inPeer)
}