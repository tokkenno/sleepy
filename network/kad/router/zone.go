package router

import (
	"com/github/reimashi/sleepy/types/uint128"
	"com/github/reimashi/sleepy/network/kad"
	"net"
	"time"
	"errors"
)

const (
	maxLevels = 6
)

// Zone is node inside a binary tree of k-buckets
type Zone struct {
	localId uint128.UInt128
	parent *Zone
	leftChild *Zone
	rightChild *Zone
	level uint
	bucket *KBucket
	randomLookupTimer *time.Ticker
}

// Load a router zone tree from file
func FromFile(localId uint128.UInt128, path string) (*Zone, error) {
	return nil, errors.New("Not implemented yet")
}

// Create a new Zone from the local peer Id
func FromId(id uint128.UInt128) *Zone {
	rz := &Zone{
		localId: id,
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

// Count the number of peers inside the branch
func (this *Zone) CountPeers() int {
	if this.isLeaf() {
		return this.bucket.Count()
	} else {
		return this.leftChild.CountPeers() + this.rightChild.CountPeers()
	}
}

// Check if the object is a leaf (is not, is a branch)
func (this *Zone) isLeaf() bool {
	return this.bucket != nil
}

// Get the max depth of this branch
func (this *Zone) maxDepth() int {
	if this.isLeaf() {
		return 0
	} else {
		ld := this.leftChild.maxDepth()
		rd := this.rightChild.maxDepth()
		if ld > rd { return ld + 1} else { return rd + 1 }
	}
}

// Run a timer to do random lookup of peers
func (this *Zone) runRandomLookupTimer() {
	for range this.randomLookupTimer.C {
		if (this.parent == nil) {
			this.onRandomLookupTimer()
		}
	}
}

// Handle the RandomLookup timer and run a lookup of a random peer inside each leaf
func (this *Zone) onRandomLookupTimer() bool {
	if this.isLeaf() {
		if (this.level < maxLevels || float32(this.bucket.Count()) >= (maxSize * 0.8)) {
			this.randomLookup();
			return true;
		} else {
			return false;
		}
	} else {
		return this.leftChild.onRandomLookupTimer() && this.rightChild.onRandomLookupTimer()
	}
}

func (this *Zone) randomLookup() {

}

func (this *Zone) onSmallTimer() {

}

// Check if the current leaf can be splitted in a branch with 2 leafs
func (this *Zone) CanSplit() bool {
	// Max levels allowed reached
	if (this.level >= 127) {
		return false
	}

	if this.level < (maxLevels - 1) && this.bucket.Count() == maxSize {
		return true
	}

	return false
}

func (this *Zone) Split() {

}

func (this *Zone) Add(peer *kad.Peer) error {
	// TODO: Filter IPs and protocol versions
	if (!this.isLeaf()) {

	} else {
		if !this.localId.Equal(peer.Id) {
			locPeer, err := this.bucket.GetPeer(&peer.Id)
			if err == nil {
				// Update
				locPeer.Update(peer)
			} else {
				if this.bucket.IsFull() {
					// Split
				} else {
					// Insert
				}
			}
		}
	}

	return nil
}

func (this *Zone) GetPeer(id uint128.UInt128) *kad.Peer {
	return nil
}

func (this *Zone) GetRandomPeer(id uint128.UInt128) *kad.Peer {
	return nil
}

// Set a peer as verified
func (this *Zone) VerifyPeer(id uint128.UInt128, ip net.IP) bool {
	peer := this.GetPeer(id)

	if (peer == nil) {
		return false
	} else if !ip.Equal(peer.GetIP()) {
		return false
	} else {
		peer.IpVerified = true
		return true
	}
}