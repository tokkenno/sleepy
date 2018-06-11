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
	localId           types.UInt128
	zoneIndex         types.UInt128
	parent            *Zone
	leftChild         *Zone
	rightChild        *Zone
	level             uint8
	bucket            *KBucket
	randomLookupTimer *time.Ticker
	randomGenerator   *rand.Rand
	checkStopFlag	  bool
}

// Load a router zone tree from file
func FromFile(localId types.UInt128, path string) (*Zone, error) {
	return nil, errors.New("not implemented yet")
}

// Create a new Zone tree (Router) from the local peer Id
func NewRouter(id types.UInt128) *Zone {
	rz := &Zone{
		localId:           id,
		zoneIndex:         *types.NewUInt128FromInt(0),
		parent:            nil,
		leftChild:         nil,
		rightChild:        nil,
		level:             0,
		bucket:            new(KBucket),
		randomLookupTimer: nil,
		randomGenerator:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	rz.startChecks()

	return rz
}

// Create a child zone from a parent instance
func newChildZone(parent Zone, isRightChild bool) *Zone {
	zoneIndexCalculated := parent.zoneIndex.Clone()
	zoneIndexCalculated.LeftShift(1)
	if isRightChild {
		zoneIndexCalculated.Add(*types.NewUInt128FromInt(1))
	}

	rz := &Zone{
		localId:           parent.localId,
		zoneIndex:         *zoneIndexCalculated,
		parent:            &parent,
		leftChild:         nil,
		rightChild:        nil,
		level:             parent.level + 1,
		bucket:            new(KBucket),
		randomLookupTimer: nil,
		randomGenerator:   parent.randomGenerator,
	}

	rz.startChecks()

	return rz
}

// Create the two child zones from the parent instance
func newChildZones(parent Zone) (*Zone, *Zone) {
	return newChildZone(parent, false), newChildZone(parent, true)
}

// Dispose resources
func (zone *Zone) Dispose() {
	if zone.isLeaf() {
		zone.stopChecks()
		zone.bucket = nil
	} else {
		zone.leftChild.Dispose()
		zone.leftChild = nil
		zone.rightChild.Dispose()
		zone.rightChild = nil
	}
}

// Start all check subroutines
func (zone *Zone) startChecks() {
	zone.checkStopFlag = false
	go zone.runRandomLookupTimer() // TODO: Como parar?
}

// Stop all check subroutines
func (zone *Zone) stopChecks() {
	zone.checkStopFlag = true
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
		if ld > rd {
			return ld + 1
		} else {
			return rd + 1
		}
	}
}

// Run a timer to do random lookup of peers
func (zone *Zone) runRandomLookupTimer() {
	zone.randomLookupTimer = time.NewTicker(10 * time.Second) // TODO: Less time in major levels?
	for range zone.randomLookupTimer.C {
		if zone.checkStopFlag {
			// TODO: Find other immediate way to stop
			break
		} else {
			zone.onRandomLookupTimer()
			zone.onSmallTimer() // TODO: ??
		}
	}
}

// Handle the RandomLookup timer and run a lookup of a random peer inside each leaf
func (zone *Zone) onRandomLookupTimer() bool {
	if zone.isLeaf() {
		if zone.level < maxLevels || float32(zone.bucket.CountPeers()) >= (maxBucketSize*0.8) {
			// TODO: Logic
			return true
		} else {
			return false
		}
	} else {
		return zone.leftChild.onRandomLookupTimer() && zone.rightChild.onRandomLookupTimer()
	}
}

func (zone *Zone) onSmallTimer() {

}

// Check if the current leaf can be splitted in a branch with 2 leafs
func (zone *Zone) canSplit() bool {
	// TODO: Check if the global contact limit is reached

	// Max levels allowed reached
	if zone.level >= 127 {
		return false
	}

	// Check if this zone is allowed to split.
	if zone.level < (maxLevels-1) && zone.bucket.CountPeers() == maxBucketSize {
		return true
	}

	return false
}

// Retract and consolidate the tree branch if it is possible
func (zone *Zone) consolidate() {
	if zone.isLeaf() {
		return
	} else {
		if zone.leftChild.isLeaf() {
			zone.leftChild.consolidate()
		}
		if zone.rightChild.isLeaf() {
			zone.rightChild.consolidate()
		}
		if zone.leftChild.isLeaf() && zone.rightChild.isLeaf() && zone.CountPeers() < (maxBucketSize/2) {
			// Initialize leaf variables
			zone.bucket = new(KBucket)

			// Stop left child, get contacts and dispose
			zone.leftChild.stopChecks()
			for _, currPeer := range zone.leftChild.bucket.peers {
				zone.bucket.AddPeer(&currPeer)
			}
			zone.leftChild.Dispose()
			zone.leftChild = nil

			// Stop right child, get contacts and dispose
			zone.rightChild.stopChecks()
			for _, currPeer := range zone.rightChild.bucket.peers {
				zone.bucket.AddPeer(&currPeer)
			}
			zone.rightChild = nil

			// Restart checks
			zone.startChecks()
		}
	}
}

// Split a leaf into a branch with two leafs
func (zone *Zone) split() error {
	if zone.canSplit() {
		zone.stopChecks()
		zone.leftChild, zone.rightChild = newChildZones(*zone)

		for _, currPeer := range zone.bucket.GetPeers() {
			distance := currPeer.GetDistance(zone.localId)
			if distance.GetBit(zone.GetLevel()) == 0 {
				zone.leftChild.AddPeer(&currPeer)
			} else {
				zone.rightChild.AddPeer(&currPeer)
			}
		}

		zone.bucket = nil
	} else {
		return errors.New("this zone can't be splitted")
	}

	return nil
}

// Get the zone depth on the router tree
func (zone *Zone) GetLevel() int {
	return int(zone.level)
}

// Get the zone max depth from the most length branch
func (zone *Zone) GetMaxDepth() int {
	if zone.isLeaf() {
		return 0
	} else {
		left := zone.leftChild.GetMaxDepth()
		right := zone.rightChild.GetMaxDepth()
		if left > right {
			return left
		} else {
			return right
		}
	}
}

func (zone *Zone) AddPeer(peer *kad.Peer) error {
	// TODO: Filter IPs and protocol versions

	if !zone.isLeaf() {
		// If the zone isn't leaf, insert in the indicated child zone
		if peer.GetDistance(zone.localId).GetBit(int(zone.level)) == 0 {
			zone.leftChild.AddPeer(peer)
		} else {
			zone.rightChild.AddPeer(peer)
		}
	} else {
		if !zone.localId.Equal(*peer.GetId()) {
			locPeer, err := zone.bucket.GetPeer(*peer.GetId())

			if err == nil && locPeer != nil {
				// If the peer already exists, update
				locPeer.Update(peer)
				return nil
			} else if !zone.bucket.IsFull() {
				// If not exists, but leaf has free space, insert
				zone.bucket.AddPeer(peer)
				return nil
			} else if zone.canSplit() {
				// If don't have free space but can be splitted, split and retry
				zone.split()
				return zone.AddPeer(peer)
			} else {
				return errors.New("the peer can't be added")
			}
		} else {
			return errors.New("the router can't contains itself")
		}
	}

	return nil
}

// Get a peer from his id
func (zone *Zone) GetPeer(id types.UInt128) (*kad.Peer, error) {
	if zone.isLeaf() {
		return zone.bucket.GetPeer(id)
	} else {
		distance := types.Xor(zone.localId, id)
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

// Get a slice of all peers
func (zone *Zone) GetAllPeers() []kad.Peer {
	if zone.isLeaf() {
		return zone.bucket.peers
	} else {
		return append(zone.leftChild.GetAllPeers(), zone.rightChild.GetAllPeers()...)
	}
}

// ¿?
func (zone *Zone) GetDepthPeers(depth int) []kad.Peer {
	if zone.isLeaf() {
		return zone.bucket.peers
	} else if depth <= 0 {
		return zone.GetRandomKBucketPeers()
	} else {
		return append(zone.leftChild.GetDepthPeers(depth-1), zone.rightChild.GetDepthPeers(depth-1)...)
	}
}

// Get the peers from a random child bucket
func (zone *Zone) GetRandomKBucketPeers() []kad.Peer {
	if zone.isLeaf() {
		return zone.bucket.peers
	} else {
		if zone.randomGenerator.Intn(1) == 0 {
			return zone.leftChild.GetRandomKBucketPeers()
		} else {
			return zone.rightChild.GetRandomKBucketPeers()
		}
	}
}

// Count the number of peers inside the branch
func (zone *Zone) CountPeers() int {
	if zone.isLeaf() {
		return zone.bucket.CountPeers()
	} else {
		return zone.leftChild.CountPeers() + zone.rightChild.CountPeers()
	}
}

// Check if exists a peer with a concrete id into the zone
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
