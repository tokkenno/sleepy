package router

import (
	"crypto/rand"
	"errors"
	"net"
	kadTypes "sleepy/network/kad/types"
	"sleepy/types"
	"sort"
	"sync"
)

const (
	maxBucketSize = 16 // Max number of peers in each k-bucket
	maxPeersPerIP = 3  // Number of peers permitted from same public IP
)

// K-bucket is a queue of k peers ordered by TTL
type kBucket struct {
	peers       []kadTypes.Peer
	peersCount  int
	peersAccess sync.Mutex
}

func newKBucket() *kBucket {
	return &kBucket{
		peers:       make([]kadTypes.Peer, maxBucketSize),
		peersCount:  0,
		peersAccess: sync.Mutex{},
	}
}

// CountPeers returns the number of peers in the bucket
func (bucket *kBucket) CountPeers() int {
	return bucket.peersCount
}

// CountRemainingPeers returns the number of peers that can be added to the bucket
func (bucket *kBucket) CountRemainingPeers() int {
	return maxBucketSize - bucket.CountPeers()
}

// AddPeer add a peer to the bucket
func (bucket *kBucket) AddPeer(newPeer kadTypes.Peer) error {
	if newPeer == nil {
		return errors.New("kBucket only can storage not null peers")
	}

	bucket.peersAccess.Lock()

	if bucket.peersCount >= maxBucketSize {
		bucket.peersAccess.Unlock()
		return errors.New("the current kBucket is full")
	}

	sameNetwork := 0
	for peerIndex := 0; peerIndex < bucket.peersCount; peerIndex++ {
		if bucket.peers[peerIndex].Equal(newPeer) {
			bucket.peersAccess.Unlock()
			return errors.New("kBucket already contains the passed peer")
		}
		if bucket.peers[peerIndex].GetIP().Equal(newPeer.GetIP()) {
			sameNetwork++
		}
	}

	if sameNetwork >= maxPeersPerIP {
		bucket.peersAccess.Unlock()
		return errors.New("many peers for the current IP")
	}

	bucket.peers[bucket.peersCount] = newPeer
	bucket.peersCount++
	bucket.peersAccess.Unlock()
	// TODO: Adjust global tracking
	return nil
}

// RemovePeer remove a peer from the bucket
func (bucket *kBucket) RemovePeer(peer kadTypes.Peer) error {
	bucket.peersAccess.Lock()

	removing := false
	for peerIndex := 0; peerIndex < bucket.peersCount; peerIndex++ {
		if removing {
			bucket.peers[peerIndex-1] = bucket.peers[peerIndex]
		} else if bucket.peers[peerIndex].Equal(peer) {
			bucket.peers[peerIndex] = nil
			removing = true
		}
	}

	if removing {
		bucket.peersCount--
		bucket.peersAccess.Unlock()
		return nil
	} else {
		bucket.peersAccess.Unlock()
		return errors.New("kBucket don't contains a peer with the passed id")
	}
}

// GetPeer find and retrieve a peer from his id
func (bucket *kBucket) GetPeer(id types.UInt128) (kadTypes.Peer, error) {
	bucket.peersAccess.Lock()

	for peerIndex := 0; peerIndex < bucket.peersCount; peerIndex++ {
		peerId := bucket.peers[peerIndex].GetID()

		if peerId.Equal(id) {
			bucket.peersAccess.Unlock()
			return bucket.peers[peerIndex], nil
		}
	}

	bucket.peersAccess.Unlock()
	return nil, errors.New("kBucket don't contains a peer with the passed id")
}

// GetPeerByAddr find and retrieve a peer from his network address
func (bucket *kBucket) GetPeerByAddr(addr net.Addr) (kadTypes.Peer, error) {
	bucket.peersAccess.Lock()

	var ip net.IP
	var port uint16
	var isTcp bool

	switch pAddr := addr.(type) {
	case *net.UDPAddr:
		ip = pAddr.IP
		port = uint16(pAddr.Port)
		isTcp = false
	case *net.TCPAddr:
		ip = pAddr.IP
		port = uint16(pAddr.Port)
		isTcp = true
	default:
		return nil, errors.New("incompatible network type")
	}

	for peerIndex := 0; peerIndex < bucket.peersCount; peerIndex++ {
		if bucket.peers[peerIndex].GetIP().Equal(ip) {
			if (isTcp && bucket.peers[peerIndex].GetTCPPort() == port) || (!isTcp && bucket.peers[peerIndex].GetUDPPort() == port) {
				bucket.peersAccess.Unlock()
				return bucket.peers[peerIndex], nil
			}
		}
	}

	bucket.peersAccess.Unlock()
	return nil, errors.New("kBucket don't contains a peer with the passed ip")
}

func (bucket *kBucket) GetRandomPeer() (kadTypes.Peer, error) {
	b := make([]byte, 1)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	bucket.peersAccess.Lock()

	if bucket.peersCount > 0 {
		peer := bucket.peers[int(b[0])%bucket.peersCount]
		bucket.peersAccess.Unlock()
		return peer, nil
	} else {
		bucket.peersAccess.Unlock()
		return nil, errors.New("kBucket don't contains any peer")
	}
}

// OldestPeer returns the oldest peer in the bucket
func (bucket *kBucket) OldestPeer() kadTypes.Peer {
	if bucket.peersCount > 0 {
		return bucket.peers[0]
	} else {
		return nil
	}
}

func (bucket *kBucket) ContainsPeer(id types.UInt128) bool {
	peer, err := bucket.GetPeer(id)
	return peer != nil && err == nil
}

func (bucket *kBucket) Peers() []kadTypes.Peer {
	peerCpy := make([]kadTypes.Peer, bucket.peersCount)
	copy(peerCpy, bucket.peers[0:bucket.peersCount])
	return peerCpy
}

// GetClosestPeers get the closest [max] peers respect the [to] id.
func (bucket *kBucket) GetClosestPeers(to types.UInt128, max int) []kadTypes.Peer {
	if bucket.peersCount > 0 {
		// Filter to leave only the active peers
		peerCpy := kadTypes.Filter(bucket.Peers(), func(peer kadTypes.Peer) bool {
			return peer.IsIPVerified() && peer.IsAlive()
		})

		// Sort by distance
		sort.Slice(peerCpy, func(i int, j int) bool {
			firstPeerId := peerCpy[i].GetID()
			secondPeerId := peerCpy[j].GetID()
			iDis := types.Xor(firstPeerId, to)
			jDis := types.Xor(secondPeerId, to)
			return iDis.Compare(jDis) < 0
		})

		// Get the [max] remaining peers
		if len(peerCpy) > max {
			return peerCpy[0:max]
		} else {
			return peerCpy
		}
	} else {
		return []kadTypes.Peer{}
	}
}

func (bucket *kBucket) IsFull() bool {
	return bucket.CountPeers() == maxBucketSize
}

// Push the peer to the end of bucket (only if already exists)
func (bucket *kBucket) pushToEnd(peer kadTypes.Peer) error {
	bucket.peersAccess.Lock()

	for position, currPeer := range bucket.peers {
		if peer.Equal(currPeer) {
			bucket.peers = append(append(bucket.peers[0:position], bucket.peers[position+1:bucket.peersCount]...), peer)
			bucket.peersAccess.Unlock()
			return nil
		}
	}

	bucket.peersAccess.Unlock()
	return errors.New("kBucket don't contains the passed peer")
}

// SetPeerAlive update last viewed status of peer, recalculate and set TTL
func (bucket *kBucket) SetPeerAlive(id types.UInt128) error {
	inPeer, err := bucket.GetPeer(id)
	if err != nil {
		return err
	}

	inPeer.UpdateType()
	return bucket.pushToEnd(inPeer)
}
