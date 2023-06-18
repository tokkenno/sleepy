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

// zone is node inside a binary tree of k-buckets
type zone struct {
	localId           types.UInt128
	zoneIndex         types.UInt128
	parent            *zone
	root              *routerImp
	leftChild         *zone
	rightChild        *zone
	level             uint8
	bucket            *kBucket
	updatePeersTimer  *time.Ticker
	randomLookupTimer *time.Ticker
	checkStopFlag     bool
	zoneAccess        sync.Mutex
}

// Create a child zone from a parent instance
func newChildZone(parent zone, isRightChild bool) *zone {
	zoneIndexCalculated := parent.zoneIndex.Clone()
	zoneIndexCalculated.LeftShift(1)
	if isRightChild {
		zoneIndexCalculated.Add(types.NewUInt128FromInt(1))
	}

	rz := &zone{
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
func newChildZones(parent zone) (*zone, *zone) {
	return newChildZone(parent, false), newChildZone(parent, true)
}

// Dispose resources
func (zn *zone) Dispose() {
	if zn.isLeaf() {
		zn.stopChecks()
		zn.bucket = nil
	} else {
		zn.leftChild.Dispose()
		zn.leftChild = nil
		zn.rightChild.Dispose()
		zn.rightChild = nil
	}
}

// Get the root zone, of type routerImp
func (zn *zone) Root() *routerImp {
	return zn.root
}

// Get the parent zone in the tree
func (zn *zone) Parent() *zone {
	return zn.parent
}

// Start all check subroutines
func (zn *zone) startChecks() {
	zn.checkStopFlag = false
	go zn.runUpdatePeersTimer()
	go zn.runRandomLookupTimer()
}

// Stop all check subroutines
func (zn *zone) stopChecks() {
	zn.checkStopFlag = true
}

// Check if the object is a leaf (is not, is a branch)
func (zn *zone) isLeaf() bool {
	return zn.bucket != nil
}

// Get the max depth of this branch
func (zn *zone) maxDepth() int {
	zn.zoneAccess.Lock()
	if zn.isLeaf() {
		zn.zoneAccess.Unlock()
		return 0
	} else {
		ld := zn.leftChild.maxDepth()
		rd := zn.rightChild.maxDepth()
		zn.zoneAccess.Unlock()
		if ld > rd {
			return ld + 1
		} else {
			return rd + 1
		}
	}
}

// Run a timer to do random lookup of peers
func (zn *zone) runRandomLookupTimer() {
	zn.randomLookupTimer = time.NewTicker(10 * time.Second)
	for range zn.randomLookupTimer.C {
		if zn.checkStopFlag {
			// TODO: Find other immediate way to stop
			break
		} else {
			zn.onRandomLookupTimer()
		}
	}
}

// Handle the RandomLookup timer and run a lookup of a random peer inside each leaf (onBigTimer)
func (zn *zone) onRandomLookupTimer() {
	zn.zoneAccess.Lock()
	if zn.isLeaf() && zn.level < maxLevels || float32(zn.bucket.CountPeers()) >= (maxBucketSize*0.8) {
		// Generate a random ID inside this zn
		randId := zn.zoneIndex.Clone()
		randId.LeftShift(uint(zn.level))
		randId.Xor(zn.localId)

		// Emit event. The KAD client will insert the peer if it finds it
		zn.Root().peerLookupRequestEvent.Emit(zn, PeerIdEventArgs{Id: randId})
	}
	zn.zoneAccess.Unlock()
}

// Run a timer to do periodic check of the peers
func (zn *zone) runUpdatePeersTimer() {
	idLow, _ := zn.localId.ToUInt64()
	zn.updatePeersTimer = time.NewTicker(time.Duration(idLow) * time.Second)
	for range zn.updatePeersTimer.C {
		if zn.checkStopFlag {
			// TODO: Find other immediate way to stop
			break
		} else {
			zn.onUpdatePeersTimer()
		}
	}
}

// Handle the UpdatePeers timer and run a check and update of the peers inside each leaf (onSmallTimer)
func (zn *zone) onUpdatePeersTimer() {
	zn.zoneAccess.Lock()
	if zn.isLeaf() {
		// Remove dead entries
		for _, peer := range zn.bucket.Peers() {
			// If peer is not alive
			if !peer.IsAlive() {
				if !peer.InUse() {
					zn.bucket.RemovePeer(peer)
				}
			} else if peer.GetExpiresAt().Equal(time.Time{}) {
				peer.SetExpiration(time.Now().Add(time.Microsecond))
			}
		}

		// Update the oldest peer
		oldestPeer := zn.bucket.OldestPeer()

		if oldestPeer != nil {
			// FIXME: pContact->GetType() == 4 ???
			if !oldestPeer.IsAlive() || oldestPeer.GetExpiresAt().After(time.Now()) {
				zn.bucket.pushToEnd(oldestPeer)
				oldestPeer = nil
			}
		}

		if oldestPeer != nil {
			oldestPeer.DegradeType()
			if oldestPeer.GetProtocolVersion() >= common.ProtocolVersion2 {
				// FIXME: The version 2 or 6 send different data, see RoutingZone.cpp:937
				router := zn.Root()
				router.peerUpdateRequestEvent.EmitSync(zn, PeerEventArgs{Peer: oldestPeer})
			}
		}
	}
	zn.zoneAccess.Unlock()
}

// Check if the current leaf can be splitted in a branch with 2 leafs
func (zn *zone) canSplit() bool {
	// TODO: Check if the global contact limit is reached

	// Max levels allowed reached
	if zn.level >= 127 {
		return false
	}

	// Check if this zn is allowed to split.
	if zn.level < (maxLevels-1) && zn.bucket.CountPeers() == maxBucketSize {
		return true
	}

	return false
}

// Retract and consolidate the tree branch if it is possible
func (zn *zone) consolidate() {
	zn.zoneAccess.Lock()
	if zn.isLeaf() {
		return
	} else {
		if zn.leftChild.isLeaf() {
			zn.leftChild.consolidate()
		}
		if zn.rightChild.isLeaf() {
			zn.rightChild.consolidate()
		}
		if zn.leftChild.isLeaf() && zn.rightChild.isLeaf() && zn.CountPeers() < (maxBucketSize/2) {
			// Initialize leaf variables
			zn.bucket = new(kBucket)

			// Stop left child, get contacts and dispose
			zn.leftChild.stopChecks()
			for _, currPeer := range zn.leftChild.bucket.peers {
				zn.bucket.AddPeer(currPeer)
			}
			zn.leftChild.Dispose()
			zn.leftChild = nil

			// Stop right child, get contacts and dispose
			zn.rightChild.stopChecks()
			for _, currPeer := range zn.rightChild.bucket.peers {
				zn.bucket.AddPeer(currPeer)
			}
			zn.rightChild = nil

			// Restart checks
			zn.startChecks()
		}
	}
	zn.zoneAccess.Unlock()
}

// Split a leaf into a branch with two leafs
func (zn *zone) split() error {
	if zn.canSplit() {
		zn.stopChecks()
		zn.leftChild, zn.rightChild = newChildZones(*zn)

		for _, currPeer := range zn.bucket.Peers() {
			distance := currPeer.GetDistance(zn.localId)
			if distance.GetBit(zn.Level()) == 0 {
				zn.leftChild.AddPeer(currPeer)
			} else {
				zn.rightChild.AddPeer(currPeer)
			}
		}

		zn.bucket = nil
	} else {
		return errors.New("this zn can't be splitted")
	}

	return nil
}

// Get the zone depth on the router tree
func (zn *zone) Level() int {
	return int(zn.level)
}

// Get the zone max depth from the most length branch
func (zn *zone) MaxDepth() int {
	if zn.isLeaf() {
		return 0
	} else {
		left := zn.leftChild.MaxDepth()
		right := zn.rightChild.MaxDepth()
		if left > right {
			return left
		} else {
			return right
		}
	}
}

// Add peer to the route table
func (zn *zone) AddPeer(peer types2.Peer) error {
	// TODO: Filter IPs and protocol versions

	if !zn.isLeaf() {
		// If the zn isn't leaf, insert in the indicated child zn
		if peer.GetDistance(zn.localId).GetBit(int(zn.level)) == 0 {
			return zn.leftChild.AddPeer(peer)
		} else {
			return zn.rightChild.AddPeer(peer)
		}
	} else {
		if !zn.localId.Equal(peer.GetID()) {
			locPeer, err := zn.bucket.GetPeer(peer.GetID())

			if err == nil && locPeer != nil {
				// If the peer already exists, update
				return locPeer.UpdateFrom(peer)
			} else if !zn.bucket.IsFull() {
				// If not exists, but leaf has free space, insert
				return zn.bucket.AddPeer(peer)
			} else if zn.canSplit() {
				// If don't have free space but can be split and retry
				zn.split()
				return zn.AddPeer(peer)
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
func (zn *zone) GetPeer(id types.UInt128) (types2.Peer, error) {
	if zn.isLeaf() {
		return zn.bucket.GetPeer(id)
	} else {
		distance := types.Xor(zn.localId, id)
		if distance.GetBit(zn.Level()) == 0 {
			return zn.leftChild.GetPeer(id)
		} else {
			return zn.rightChild.GetPeer(id)
		}
	}
}

// Get a peer from his addr
func (zn *zone) GetPeerByAddr(addr net.Addr) (types2.Peer, error) {
	if zn.isLeaf() {
		return zn.bucket.GetPeerByAddr(addr)
	} else {
		left, err := zn.leftChild.GetPeerByAddr(addr)

		if err != nil {
			return zn.rightChild.GetPeerByAddr(addr)
		} else {
			return left, err
		}
	}
}

// Get a random peer from a random branch
func (zn *zone) GetRandomPeer() (types2.Peer, error) {
	if zn.isLeaf() {
		return zn.bucket.GetRandomPeer()
	} else {
		childs := [2]*zone{zn.leftChild, zn.rightChild}
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
func (zn *zone) Peers() []types2.Peer {
	if zn.isLeaf() {
		return zn.bucket.peers
	} else {
		return append(zn.leftChild.Peers(), zn.rightChild.Peers()...)
	}
}

// Get a slice with the [maxPeers] top peers.
func (zn *zone) GetTopPeers(maxPeers int, maxDepth int) []types2.Peer {
	var peers []types2.Peer

	zn.zoneAccess.Lock()
	if zn.isLeaf() {
		peers = zn.bucket.Peers()
	} else if maxDepth <= 0 {
		peers = zn.GetRandomBucketPeers()
	} else {
		fmt.Println(zn.leftChild)
		peers = zn.leftChild.GetTopPeers(maxPeers, maxDepth-1)

		if len(peers) < maxPeers {
			peers = append(peers, zn.rightChild.GetTopPeers(maxPeers-len(peers), maxDepth-1)...)
		}
	}
	zn.zoneAccess.Unlock()

	if len(peers) < maxPeers {
		return peers
	} else {
		return peers[:maxPeers]
	}
}

// Obtain peers from a random bucket located at least at the [depth] indicated on the zone tree
func (zn *zone) GetDepthPeers(depth int) []types2.Peer {
	if zn.isLeaf() {
		return zn.bucket.Peers()
	} else if depth <= 0 {
		return zn.GetRandomBucketPeers()
	} else {
		return append(zn.leftChild.GetDepthPeers(depth-1), zn.rightChild.GetDepthPeers(depth-1)...)
	}
}

// Get the peers from a random child bucket
func (zn *zone) GetRandomBucketPeers() []types2.Peer {
	if zn.isLeaf() {
		return zn.bucket.Peers()
	} else {
		if zn.Root().randomGenerator.Intn(1) == 0 {
			return zn.leftChild.GetRandomBucketPeers()
		} else {
			return zn.rightChild.GetRandomBucketPeers()
		}
	}
}

// Get the closest [max] peers respect the [to] id.
func (zn *zone) GetClosestPeers(to types.UInt128, max int) []types2.Peer {
	if zn.isLeaf() {
		// If leaf, get from bucket
		return zn.bucket.GetClosestPeers(to, max)
	} else {
		children := [2]*zone{zn.leftChild, zn.rightChild}
		rPos := to.GetBit(int(zn.level))

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
func (zn *zone) CountPeers() int {
	if zn.isLeaf() {
		return zn.bucket.CountPeers()
	} else {
		return zn.leftChild.CountPeers() + zn.rightChild.CountPeers()
	}
}

// Check if exists a peer with a concrete id into the zone
func (zn *zone) ContainsPeer(id types.UInt128) bool {
	if zn.isLeaf() {
		return zn.bucket.ContainsPeer(id)
	} else {
		return zn.leftChild.ContainsPeer(id) || zn.rightChild.ContainsPeer(id)
	}
}

// Set a peer as verified
func (zn *zone) VerifyPeer(id types.UInt128, ip net.IP) bool {
	peer, err := zn.GetPeer(id)

	if err != nil {
		return false
	} else {
		return peer.VerifyIp(ip)
	}
}
