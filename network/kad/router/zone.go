package router

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sleepy/network/ed2k/common"
	types2 "sleepy/network/kad/types"
	"sleepy/types"
	"sleepy/utils/event"
	"sync"
	"time"
)

const (
	maxLevels = 6
)

type PeerEventArgs struct {
	event.Args
	Peer types2.Peer
}

type PeerIdEventArgs struct {
	event.Args
	Id types.UInt128
}

// Zone is node inside a binary tree of k-buckets
type Zone struct {
	localId           types.UInt128
	zoneIndex         types.UInt128
	parent            *Zone
	root              *Router
	leftChild         *Zone
	rightChild        *Zone
	level             uint8
	bucket            *kBucket
	updatePeersTimer  *time.Ticker
	randomLookupTimer *time.Ticker
	checkStopFlag     bool
	zoneAccess        sync.Mutex
}

// Create a child zone from a parent instance
func newChildZone(parent Zone, isRightChild bool) *Zone {
	zoneIndexCalculated := parent.zoneIndex.Clone()
	zoneIndexCalculated.LeftShift(1)
	if isRightChild {
		zoneIndexCalculated.Add(types.NewUInt128FromInt(1))
	}

	rz := &Zone{
		localId:           parent.localId,
		zoneIndex:         zoneIndexCalculated,
		parent:            &parent,
		root:              parent.Root(),
		leftChild:         nil,
		rightChild:        nil,
		level:             parent.level + 1,
		bucket:            newKBucket(),
		updatePeersTimer:  nil,
		randomLookupTimer: nil,
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

// Get the root Zone, of type Router
func (zone *Zone) Root() *Router {
	return zone.root
}

// Get the parent Zone in the tree
func (zone *Zone) Parent() *Zone {
	return zone.parent
}

// Start all check subroutines
func (zone *Zone) startChecks() {
	zone.checkStopFlag = false
	go zone.runUpdatePeersTimer()
	go zone.runRandomLookupTimer()
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
	zone.zoneAccess.Lock()
	if zone.isLeaf() {
		zone.zoneAccess.Unlock()
		return 0
	} else {
		ld := zone.leftChild.maxDepth()
		rd := zone.rightChild.maxDepth()
		zone.zoneAccess.Unlock()
		if ld > rd {
			return ld + 1
		} else {
			return rd + 1
		}
	}
}

// Run a timer to do random lookup of peers
func (zone *Zone) runRandomLookupTimer() {
	zone.randomLookupTimer = time.NewTicker(10 * time.Second)
	for range zone.randomLookupTimer.C {
		if zone.checkStopFlag {
			// TODO: Find other immediate way to stop
			break
		} else {
			zone.onRandomLookupTimer()
		}
	}
}

// Handle the RandomLookup timer and run a lookup of a random peer inside each leaf (onBigTimer)
func (zone *Zone) onRandomLookupTimer() {
	zone.zoneAccess.Lock()
	if zone.isLeaf() && zone.level < maxLevels || float32(zone.bucket.CountPeers()) >= (maxBucketSize*0.8) {
		// Generate a random ID inside this zone
		randId := zone.zoneIndex.Clone()
		randId.LeftShift(uint(zone.level))
		randId.Xor(zone.localId)

		// Emit event. The KAD client will insert the peer if it finds it
		zone.Root().peerLookupRequestEvent.Emit(zone, PeerIdEventArgs{Id: randId})
	}
	zone.zoneAccess.Unlock()
}

// Run a timer to do periodic check of the peers
func (zone *Zone) runUpdatePeersTimer() {
	idLow, _ := zone.localId.ToUInt64()
	zone.updatePeersTimer = time.NewTicker(time.Duration(idLow) * time.Second)
	for range zone.updatePeersTimer.C {
		if zone.checkStopFlag {
			// TODO: Find other immediate way to stop
			break
		} else {
			zone.onUpdatePeersTimer()
		}
	}
}

// Handle the UpdatePeers timer and run a check and update of the peers inside each leaf (onSmallTimer)
func (zone *Zone) onUpdatePeersTimer() {
	zone.zoneAccess.Lock()
	if zone.isLeaf() {
		// Remove dead entries
		for _, peer := range zone.bucket.Peers() {
			// If peer is not alive
			if !peer.IsAlive() {
				if !peer.InUse() {
					zone.bucket.RemovePeer(peer)
				}
			} else if peer.GetExpiresAt().Equal(time.Time{}) {
				peer.SetExpiration(time.Now().Add(time.Microsecond))
			}
		}

		// Update the oldest peer
		oldestPeer := zone.bucket.OldestPeer()

		if oldestPeer != nil {
			// FIXME: pContact->GetType() == 4 ???
			if !oldestPeer.IsAlive() || oldestPeer.GetExpiresAt().After(time.Now()) {
				zone.bucket.pushToEnd(oldestPeer)
				oldestPeer = nil
			}
		}

		if oldestPeer != nil {
			oldestPeer.DegradeType()
			if oldestPeer.GetProtocolVersion() >= common.ProtocolVersion2 {
				// FIXME: The version 2 or 6 send different data, see RoutingZone.cpp:937
				router := zone.Root()
				router.peerUpdateRequestEvent.EmitSync(zone, PeerEventArgs{Peer: oldestPeer})
			}
		}
	}
	zone.zoneAccess.Unlock()
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
	zone.zoneAccess.Lock()
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
			zone.bucket = new(kBucket)

			// Stop left child, get contacts and dispose
			zone.leftChild.stopChecks()
			for _, currPeer := range zone.leftChild.bucket.peers {
				zone.bucket.AddPeer(currPeer)
			}
			zone.leftChild.Dispose()
			zone.leftChild = nil

			// Stop right child, get contacts and dispose
			zone.rightChild.stopChecks()
			for _, currPeer := range zone.rightChild.bucket.peers {
				zone.bucket.AddPeer(currPeer)
			}
			zone.rightChild = nil

			// Restart checks
			zone.startChecks()
		}
	}
	zone.zoneAccess.Unlock()
}

// Split a leaf into a branch with two leafs
func (zone *Zone) split() error {
	if zone.canSplit() {
		zone.stopChecks()
		zone.leftChild, zone.rightChild = newChildZones(*zone)

		for _, currPeer := range zone.bucket.Peers() {
			distance := currPeer.GetDistance(zone.localId)
			if distance.GetBit(zone.Level()) == 0 {
				zone.leftChild.AddPeer(currPeer)
			} else {
				zone.rightChild.AddPeer(currPeer)
			}
		}

		zone.bucket = nil
	} else {
		return errors.New("this zone can't be splitted")
	}

	return nil
}

// Get the zone depth on the router tree
func (zone *Zone) Level() int {
	return int(zone.level)
}

// Get the zone max depth from the most length branch
func (zone *Zone) MaxDepth() int {
	if zone.isLeaf() {
		return 0
	} else {
		left := zone.leftChild.MaxDepth()
		right := zone.rightChild.MaxDepth()
		if left > right {
			return left
		} else {
			return right
		}
	}
}

// Add peer to the route table
func (zone *Zone) AddPeer(peer types2.Peer) error {
	// TODO: Filter IPs and protocol versions

	if !zone.isLeaf() {
		// If the zone isn't leaf, insert in the indicated child zone
		if peer.GetDistance(zone.localId).GetBit(int(zone.level)) == 0 {
			return zone.leftChild.AddPeer(peer)
		} else {
			return zone.rightChild.AddPeer(peer)
		}
	} else {
		if !zone.localId.Equal(peer.Id()) {
			locPeer, err := zone.bucket.GetPeer(peer.Id())

			if err == nil && locPeer != nil {
				// If the peer already exists, update
				return locPeer.UpdateFrom(peer)
			} else if !zone.bucket.IsFull() {
				// If not exists, but leaf has free space, insert
				return zone.bucket.AddPeer(peer)
			} else if zone.canSplit() {
				// If don't have free space but can be split and retry
				zone.split()
				return zone.AddPeer(peer)
			} else {
				return errors.New("the peer can't be added. bucket full and can't be split")
			}
		} else {
			return errors.New("the router can't contains itself")
		}
	}

	return nil
}

// Get a peer from his id
func (zone *Zone) GetPeer(id types.UInt128) (types2.Peer, error) {
	if zone.isLeaf() {
		return zone.bucket.GetPeer(id)
	} else {
		distance := types.Xor(zone.localId, id)
		if distance.GetBit(zone.Level()) == 0 {
			return zone.leftChild.GetPeer(id)
		} else {
			return zone.rightChild.GetPeer(id)
		}
	}
}

// Get a peer from his addr
func (zone *Zone) GetPeerByAddr(addr net.Addr) (types2.Peer, error) {
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
func (zone *Zone) GetRandomPeer() (types2.Peer, error) {
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
func (zone *Zone) Peers() []types2.Peer {
	if zone.isLeaf() {
		return zone.bucket.peers
	} else {
		return append(zone.leftChild.Peers(), zone.rightChild.Peers()...)
	}
}

// Get a slice with the [maxPeers] top peers.
func (zone *Zone) GetTopPeers(maxPeers int, maxDepth int) []types2.Peer {
	var peers []types2.Peer

	zone.zoneAccess.Lock()
	if zone.isLeaf() {
		peers = zone.bucket.Peers()
	} else if maxDepth <= 0 {
		peers = zone.GetRandomBucketPeers()
	} else {
		fmt.Println(zone.leftChild)
		peers = zone.leftChild.GetTopPeers(maxPeers, maxDepth-1)

		if len(peers) < maxPeers {
			peers = append(peers, zone.rightChild.GetTopPeers(maxPeers-len(peers), maxDepth-1)...)
		}
	}
	zone.zoneAccess.Unlock()

	if len(peers) < maxPeers {
		return peers
	} else {
		return peers[:maxPeers]
	}
}

// Obtain peers from a random bucket located at least at the [depth] indicated on the Zone tree
func (zone *Zone) GetDepthPeers(depth int) []types2.Peer {
	if zone.isLeaf() {
		return zone.bucket.Peers()
	} else if depth <= 0 {
		return zone.GetRandomBucketPeers()
	} else {
		return append(zone.leftChild.GetDepthPeers(depth-1), zone.rightChild.GetDepthPeers(depth-1)...)
	}
}

// Get the peers from a random child bucket
func (zone *Zone) GetRandomBucketPeers() []types2.Peer {
	if zone.isLeaf() {
		return zone.bucket.Peers()
	} else {
		if zone.Root().randomGenerator.Intn(1) == 0 {
			return zone.leftChild.GetRandomBucketPeers()
		} else {
			return zone.rightChild.GetRandomBucketPeers()
		}
	}
}

// Get the closest [max] peers respect the [to] id.
func (zone *Zone) GetClosestPeers(to types.UInt128, max int) []types2.Peer {
	if zone.isLeaf() {
		// If leaf, get from bucket
		return zone.bucket.GetClosestPeers(to, max)
	} else {
		children := [2]*Zone{zone.leftChild, zone.rightChild}
		rPos := to.GetBit(int(zone.level))

		// Get from the closest branch
		peers := children[rPos].GetClosestPeers(to, max)
		if len(peers) < max {
			// If it has insufficient, get any from the other branch
			return append(peers, children[rPos^1].GetClosestPeers(to, max-len(peers))...)
		} else {
			return peers
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
