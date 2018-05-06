package router

import (
	"com/github/reimashi/sleepy/types"
	"com/github/reimashi/sleepy/network/kad"
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
func (this *Zone) newChilds() (*Zone, *Zone) {
	leftChild := &Zone {
		localId: this.localId,
		maskId: nil, // TODO: Calculate
		parent: this,
		leftChild: nil,
		rightChild: nil,
		level: this.level + 1,
		bucket: new(KBucket),
		randomLookupTimer: time.NewTicker(10 * time.Second),
	}
	go leftChild.runRandomLookupTimer() // TODO: Como parar?

	rightChild := &Zone {
		localId: this.localId,
		maskId: nil, // TODO: Calculate
		parent: this,
		leftChild: nil,
		rightChild: nil,
		level: this.level + 1,
		bucket: new(KBucket),
		randomLookupTimer: time.NewTicker(10 * time.Second),
	}
	go rightChild.runRandomLookupTimer() // TODO: Como parar?

	return leftChild, rightChild
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

// Split a leaf into a branch with two leafs
func (this *Zone) Split() error {
	if this.CanSplit() {
		// TODO: Stop timer
		this.leftChild, this.rightChild = this.newChilds()

		for _, currPeer := range this.bucket.GetPeers() {
			if currPeer.Distance(this.localId).GetBit(this.level) == 0 {
				this.leftChild.AddPeer(&currPeer)
			} else {
				this.rightChild.AddPeer(&currPeer)
			}
		}

		this.bucket = nil
	} else {
		return errors.New("This zone can't be splitted")
	}

	return nil
}

func (this *Zone) AddPeer(peer *kad.Peer) error {
	// TODO: Filter IPs and protocol versions
	if (!this.isLeaf()) {

	} else {
		if !this.localId.Equal(peer.Id) {
			locPeer, err := this.bucket.GetPeer(peer.Id)
			if err == nil {
				// Update
				locPeer.Update(peer)
			} else {
				if this.bucket.IsFull() {
					// Split
					this.Split()
				} else {
					// Insert
				}
			}
		}
	}

	return nil
}

// Get a peer from his id
func (this *Zone) GetPeer(id types.UInt128) (*kad.Peer, error) {
	if this.isLeaf() {
		return this.bucket.GetPeer(id)
	} else {
		distance := types.Xor(this.maskId, id) // TODO: localId?
		if distance.GetBit(this.level) == 0 {
			return this.leftChild.GetPeer(id)
		} else {
			return this.rightChild.GetPeer(id)
		}
	}
}

// Get a peer from his ip
func (this *Zone) GetPeerByIp(ip net.IP) (*kad.Peer, error) {
	if this.isLeaf() {
		return this.bucket.GetPeerByIp(ip)
	} else {
		left, err := this.leftChild.GetPeerByIp(ip)

		if err != nil {
			return this.rightChild.GetPeerByIp(ip)
		} else {
			return left, err
		}
	}
}

// Get a random peer from a random branch
func (this *Zone) GetRandomPeer() (*kad.Peer, error) {
	if this.isLeaf() {
		return this.bucket.GetRandomPeer()
	} else {
		childs := [2]*Zone{this.leftChild, this.rightChild}
		rPos := rand.Intn(1)

		peer, err := childs[rPos].GetRandomPeer()
		if err != nil {
			return childs[1^rPos].GetRandomPeer()
		} else {
			return peer, err
		}
	}
}

// Set a peer as verified
func (this *Zone) VerifyPeer(id types.UInt128, ip net.IP) bool {
	peer, err := this.GetPeer(id)

	if (err != nil) {
		return false
	} else if !ip.Equal(peer.GetIP()) {
		return false
	} else {
		peer.IpVerified = true
		return true
	}
}