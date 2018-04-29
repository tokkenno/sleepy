package router

import (
	"com/github/reimashi/sleepy/types/uint128"
	"com/github/reimashi/sleepy/network/kad"
	"net"
	"time"
)

const (
	maxLevels = 6
	maxSize = 16
)

type RouterZone struct {
	localId uint128.UInt128
	parent *RouterZone
	leftChild *RouterZone
	rightChild *RouterZone
	level uint
	bucket *KBucket
	randomLookupTimer *time.Ticker
}

func FromFile(localId uint128.UInt128, path string) *RouterZone {
	return nil
}

func FromId(id uint128.UInt128) *RouterZone {
	rz := &RouterZone{
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

func (this *RouterZone) GetSize() int {
	return 0
}

func (this *RouterZone) CountPeers() int {
	if this.isLeaf() {
		return this.bucket.CountPeers()
	} else {
		return this.leftChild.CountPeers() + this.rightChild.CountPeers()
	}
}

func (this *RouterZone) isLeaf() bool {
	return this.bucket != nil
}

func (this *RouterZone) maxDepth() int {
	if this.isLeaf() {
		return 0
	} else {
		ld := this.leftChild.maxDepth()
		rd := this.rightChild.maxDepth()
		if ld > rd { return ld + 1} else { return rd + 1 }
	}
}

func (this *RouterZone) runRandomLookupTimer() {
	for range this.randomLookupTimer.C {
		if (this.parent == nil) {
			this.onRandomLookupTimer()
		}
	}
}

func (this *RouterZone) onRandomLookupTimer() bool {
	if this.isLeaf() {
		if (this.level < maxLevels || float32(this.bucket.CountPeers()) >= (maxSize * 0.8)) {
			this.randomLookup();
			return true;
		} else {
			return false;
		}
	} else {
		return this.leftChild.onRandomLookupTimer() && this.rightChild.onRandomLookupTimer()
	}
}

func (this *RouterZone) randomLookup() {

}

func (this *RouterZone) onSmallTimer() {

}

func (this *RouterZone) CanSplit() bool {
	// Max levels allowed reached
	if (this.level >= 127) {
		return false
	}

	if this.level < maxLevels - 1 && this.GetSize() == maxSize {
		return true
	}

	return false
}

func (this *RouterZone) Split() {

}

func (this *RouterZone) Add(peer *kad.Peer) error {
	return nil
}

func (this *RouterZone) GetPeer(id uint128.UInt128) *kad.Peer {
	return nil
}

func (this *RouterZone) GetRandomPeer(id uint128.UInt128) *kad.Peer {
	return nil
}

func (this *RouterZone) VerifyPeer(id uint128.UInt128, ip net.IP) bool {
	peer := this.GetPeer(id)

	if (peer == nil) {
		return false
	} else if !ip.Equal(peer.UDPAddr.IP) {
		return false
	} else {
		peer.IpVerified = true
		return true
	}
}