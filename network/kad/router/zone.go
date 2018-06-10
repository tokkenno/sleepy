package router

import (
	"github.com/tokkenno/sleepy/types"
	"github.com/tokkenno/sleepy/network/kad"
	"net"
	"time"
	"errors"
	"math/rand"
)

const (
	maxLevels = 6
)

// Zone is node inside a binary tree of k-buckets
type Zone struct {
	localId types.UInt128
	maskId types.UInt128
	parent *Zone
	leftChild *Zone
	rightChild *Zone
	level uint8
	bucket *KBucket
	randomLookupTimer *time.Ticker
}

// Load a router zone tree from file
func FromFile(localId types.UInt128, path string) (*Zone, error) {
	return nil, errors.New("Not implemented yet")
}

// Create a new Zone from the local peer Id
func FromId(id types.UInt128) *Zone {
	rz := &Zone{
		localId: id,
		maskId: id,
		parent: nil,
		leftChild: nil,
		rightChild: nil,
		level: 0,
		bucket: new(KBucket),
		randomLookupTimer: time.NewTicker(10 * time.Second),
	}
	go rz.runRandomLookupTimer() // TODO: Como parar?
	return rz
}

// Create a new child zone
func (zone *Zone) newChilds() (*Zone, *Zone) {
	leftChild := &Zone {
		localId: zone.localId,
		//maskId: nil, // TODO: Calculate
		parent:            zone,
		leftChild:         nil,
		rightChild:        nil,
		level:             zone.level + 1,
		bucket:            new(KBucket),
		randomLookupTimer: time.NewTicker(10 * time.Second),
	}
	go leftChild.runRandomLookupTimer() // TODO: Como parar?

	rightChild := &Zone {
		localId: zone.localId,
		//maskId: nil, // TODO: Calculate
		parent:            zone,
		leftChild:         nil,
		rightChild:        nil,
		level:             zone.level + 1,
		bucket:            new(KBucket),
		randomLookupTimer: time.NewTicker(10 * time.Second),
	}
	go rightChild.runRandomLookupTimer() // TODO: Como parar?

	return leftChild, rightChild
}

// Count the number of peers inside the branch
func (zone *Zone) CountPeers() int {
	if zone.isLeaf() {
		return zone.bucket.CountPeers()
	} else {
		return zone.leftChild.CountPeers() + zone.rightChild.CountPeers()
	}
}

// Check if the object is a leaf (is not, is a branch)
func (zone *Zone) isLeaf() bool {
	return zone.bucket != nil
}

// Get the max depth of this branch
func (zone *Zone) maxDepth() int {
	if zone.isLeaf() {
		return 0
	} else {
		ld := zone.leftChild.maxDepth()
		rd := zone.rightChild.maxDepth()
		if ld > rd { return ld + 1} else { return rd + 1 }
	}
}

// Get the zone depth on the router tree
func (zone *Zone) GetLevel() int {
	return int(zone.level)
}

// Run a timer to do random lookup of peers
func (zone *Zone) runRandomLookupTimer() {
	for range zone.randomLookupTimer.C {
		if zone.parent == nil {
			zone.onRandomLookupTimer()
		}
	}
}

// Handle the RandomLookup timer and run a lookup of a random peer inside each leaf
func (zone *Zone) onRandomLookupTimer() bool {
	if zone.isLeaf() {
		if zone.level < maxLevels || float32(zone.bucket.CountPeers()) >= (maxSize * 0.8) {
			zone.randomLookup()
			return true
		} else {
			return false
		}
	} else {
		return zone.leftChild.onRandomLookupTimer() && zone.rightChild.onRandomLookupTimer()
	}
}

func (zone *Zone) randomLookup() {

}

func (zone *Zone) onSmallTimer() {

}

// Check if the current leaf can be splitted in a branch with 2 leafs
func (zone *Zone) CanSplit() bool {
	// Max levels allowed reached
	if zone.level >= 127 {
		return false
	}

	if zone.level < (maxLevels - 1) && zone.bucket.CountPeers() == maxSize {
		return true
	}

	return false
}

// Split a leaf into a branch with two leafs
func (zone *Zone) Split() error {
	if zone.CanSplit() {
		// TODO: Stop timer
		zone.leftChild, zone.rightChild = zone.newChilds()

		for _, currPeer := range zone.bucket.GetPeers() {
			distance := currPeer.Distance(zone.localId)
			if distance.GetBit(zone.GetLevel()) == 0 {
				zone.leftChild.AddPeer(&currPeer)
			} else {
				zone.rightChild.AddPeer(&currPeer)
			}
		}

		zone.bucket = nil
	} else {
		return errors.New("This zone can't be splitted")
	}

	return nil
}

func (zone *Zone) AddPeer(peer *kad.Peer) error {
	// TODO: Filter IPs and protocol versions
	if !zone.isLeaf() {

	} else {
		if !zone.localId.Equal(*peer.GetId()) {
			locPeer, err := zone.bucket.GetPeer(*peer.GetId())
			if err == nil {
				// Update
				locPeer.Update(peer)
			} else {
				if zone.bucket.IsFull() {
					// Split
					zone.Split()
				} else {
					// Insert
					zone.bucket.AddPeer(peer)
				}
			}
		}
	}

	return nil
}

// Get a peer from his id
func (zone *Zone) GetPeer(id types.UInt128) (*kad.Peer, error) {
	if zone.isLeaf() {
		return zone.bucket.GetPeer(id)
	} else {
		distance := types.Xor(zone.maskId, id) // TODO: localId?
		if distance.GetBit(zone.GetLevel()) == 0 {
			return zone.leftChild.GetPeer(id)
		} else {
			return zone.rightChild.GetPeer(id)
		}
	}
}

// Get a peer from his addr
func (zone *Zone) GetPeerByAddr(addr net.Addr) (*kad.Peer, error) {
	if zone.isLeaf() {
		return zone.bucket.GetPeerByAddr(addr)
	} else {
		left, err := zone.leftChild.GetPeerByAddr(addr)

		if err != nil {
			return zone.rightChild.GetPeerByAddr(addr)
		} else {
			return left, err
		}
	}
}

// Get a random peer from a random branch
func (zone *Zone) GetRandomPeer() (*kad.Peer, error) {
	if zone.isLeaf() {
		return zone.bucket.GetRandomPeer()
	} else {
		childs := [2]*Zone{zone.leftChild, zone.rightChild}
		rPos := rand.Intn(1)

		peer, err := childs[rPos].GetRandomPeer()
		if err != nil {
			return childs[1^rPos].GetRandomPeer()
		} else {
			return peer, err
		}
	}
}

func (zone *Zone) ContainsPeer(id types.UInt128) bool {
	if zone.isLeaf() {
		return zone.bucket.ContainsPeer(id)
	} else {
		return zone.leftChild.ContainsPeer(id) || zone.rightChild.ContainsPeer(id)
	}
}

// Set a peer as verified
func (zone *Zone) VerifyPeer(id types.UInt128, ip net.IP) bool {
	peer, err := zone.GetPeer(id)

	if err != nil {
		return false
	} else {
		return peer.VerifyIp(ip)
	}
}